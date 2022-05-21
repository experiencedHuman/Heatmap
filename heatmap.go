package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/kvogli/Heatmap/proto/api"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kvogli/Heatmap/DBService"
	"github.com/kvogli/Heatmap/NavigaTUM"
	"github.com/kvogli/Heatmap/RoomFinder"
)


func getDataFromURL(filename, url string) {
	resp, httpError := http.Get(url)
	if httpError != nil {
		log.Fatalf("Could not retrieve csv data from URL! %q", httpError)
	}

	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)
	csvReader.Comma = ','

	file, osError := os.Create(filename)
	if osError != nil {
		log.Fatalf("Could not create file, err: %q", osError)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Comma = ';'
	defer csvWriter.Flush()

	var csvRecord []string
	for i := 0; ; i++ {
		csvRecord, httpError = csvReader.Read()
		if httpError == io.EOF {
			break
		} else if httpError != nil {
			panic(httpError)
		} else {
			fields := strings.Fields(csvRecord[0]) // get substrings separated by whitespaces
			network := fields[0]
			current := strings.Split(fields[1], ":")[1]
			max := strings.Split(fields[2], ":")[1]
			min := strings.Split(fields[3], ":")[1]

			dateAndTime := csvRecord[1]
			avg := csvRecord[2]

			csvWriter.Write([]string{
				network, current, max, min, dateAndTime, avg,
			})
		}
	}
}

type JsonEntry struct {
	Intensity float64
	Lat       float64
	Long      float64
	Floor     string
}

const (
	ApstatTable = "./data/sqlite/apstat.db"
	ApstatCSV   = "data/csv/apstat.csv"
)

var apstatDB = DBService.InitDB(ApstatTable)

func saveAPsToJSON(dst string, totalLoad int) {
	APs := DBService.RetrieveAPs(apstatDB, true)
	var jsonData []JsonEntry

	for _, ap := range APs {
		currTotLoad := RoomFinder.GetCurrentTotalLoad(ap.Load)
		var intensity float64 = float64(currTotLoad) / float64(totalLoad)
		lat, _ := strconv.ParseFloat(ap.Lat, 64)
		lng, _ := strconv.ParseFloat(ap.Long, 64)
		jsonEntry := JsonEntry{intensity, lat, lng, ap.Floor}
		jsonData = append(jsonData, jsonEntry)
	}

	bytes, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(dst, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func scrapeRoomFinder() (RoomFinder.Result, int) {
	db := DBService.InitDB(ApstatTable)
	APs := DBService.RetrieveAPs(db, false)
	roomInfos, totalLoad := RoomFinder.PrepareDataToScrape(APs)
	res := RoomFinder.ScrapeURLs(roomInfos)
	
	log.Println("Number of retrieved APs:", len(APs))
	log.Println("Number of retrieved URLs:", len(res.Successes))
	
	for _, val := range res.Successes {
		where := fmt.Sprintf("ID='%s'", val.ID)
		DBService.UpdateColumn(db, "apstat", "Lat", val.Lat, where)
		DBService.UpdateColumn(db, "apstat", "Long", val.Long, where)
	}

	return res, totalLoad
}

func scrapeNavigaTUM(res RoomFinder.Result) (count int) {
	count = 0 // number of found coordinates
	db := DBService.InitDB(ApstatTable)
	
	for _, res := range res.Failures {
		var roomID string
		if strings.Contains(res.RoomNr, "OG") || res.RoomNr == "" || strings.Contains(res.RoomNr, ".."){
			roomID = res.BuildingNr
		} else {
			roomID = fmt.Sprintf("%s.%s", res.BuildingNr, res.RoomNr)
		}
		
		lat, long, found := NavigaTUM.GetRoomCoordinates(roomID)

		if found {
			where := fmt.Sprintf("ID='%s'", res.ID)
			DBService.UpdateColumn(db, "apstat", "Lat", lat, where)
			DBService.UpdateColumn(db, "apstat", "Long", long, where)
			count++
		} else {
			lat, long, found = NavigaTUM.GetRoomCoordinates(res.BuildingNr)
			if found {
				where := fmt.Sprintf("ID='%s'", res.ID)
				DBService.UpdateColumn(db, "apstat", "Lat", lat, where)
				DBService.UpdateColumn(db, "apstat", "Long", long, where)
				count++
			}
		}
	}

	return 
}


type server struct {
	pb.UnimplementedAPServiceServer
}

func NewServer() *server {
	return &server{}
}

// TODO implement request with ID of access point and appropriate response
func (s *server) GetAccessPoint(ctx context.Context, in *pb.APRequest) (*pb.AccessPoint, error) {
	log.Printf("Received request from client for AP with id: %s", in.Id)
	return &pb.AccessPoint{Name: "apa99-k99"}, nil
}

func (s *server) ListAccessPoints(in *emptypb.Empty, stream pb.APService_ListAccessPointsServer) error {
	db := DBService.InitDB(ApstatTable)
	apList := DBService.RetrieveAPs(db, true)
	log.Printf("Sending %d APs ...", len(apList))
	i := 1
	for _, ap := range apList {
		nty := fmt.Sprintf("%d", i)
		i++
		if err:= stream.Send(
			&pb.APResponse{
				Accesspoint: 
					&pb.AccessPoint{
						Name: ap.Name, 
						Lat: ap.Lat, 
						Long: ap.Long, 
						Intensity: nty},
			}); err != nil { //TODO implement intensity
			return err
		}
	}
	return nil
}


func main() {
	port := 50051
	
	fmt.Println("Starting server...")
	
	lis, err := net.Listen("tcp", fmt.Sprintf("192.168.0.109:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAPServiceServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())

	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"192.168.0.109:50051",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterAPServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	gwServer := &http.Server {
		Addr: "192.168.0.109:50052",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://192.168.0.109:50052")
	log.Fatalln(gwServer.ListenAndServe())

	// result, totalLoad := scrapeRoomFinder() //Note that room finder must first be scraped to jump to navigatum this way
	// cnt := scrapeNavigaTUM(result)
	// fmt.Println(cnt)
	
	// saveAPsToJSON("data/json/ap.json", totalLoad)
}
