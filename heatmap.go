package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/kvogli/Heatmap/proto/api"

	// "google.golang.org/protobuf/types/known/emptypb"

	"github.com/kvogli/Heatmap/DBService"
	"github.com/kvogli/Heatmap/LRZscraper"
	// "github.com/kvogli/Heatmap/RoomFinder"
)

type JsonEntry struct {
	Intensity float64
	Lat       float64
	Long      float64
	Floor     string
}

const (
	heatmapDB = "./data/sqlite/heatmap.db"
)

// Retrieves all access points from the database
// and stores them in JSON format in `dst` e.g. "data/json/ap.json"
func saveAPsToJSON(dst string, totalLoad int) {
	APs := DBService.RetrieveAPsOfTUM(true)
	var jsonData []JsonEntry

	for _, ap := range APs {
		currTotLoad := 0
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

type server struct {
	pb.UnimplementedAPServiceServer
}

func (s *server) GetAccessPoint(ctx context.Context, in *pb.APRequest) (*pb.AccessPoint, error) {
	name := in.Name
	ts := in.Timestamp
	log.Printf("Received request for AP with name: %s and timestamp: %s", name, in.Timestamp)

	db := DBService.InitDB(heatmapDB)

	day, hr := getDayAndHourFromTimestamp(ts)
	ap := DBService.GetHistoryForSingleAP(name, day, hr)
	db.Close()

	load, _ := strconv.Atoi(ap.Load)

	return &pb.AccessPoint{
		Name:      ap.Name,
		Lat:       ap.Lat,
		Long:      ap.Long,
		Intensity: int64(load),
		Max:       int64(ap.Max),
		Min:       int64(ap.Min),
	}, nil
}

func getDayAndHourFromTimestamp(timestamp string) (int, int) {
	ts := strings.Split(timestamp, " ")
	date := ts[0]
	yearMonthDay := strings.Split(date, "-")
	day, err := strconv.Atoi(yearMonthDay[2])
	if err != nil {
		day = 0
	}
	hr, err := strconv.Atoi(ts[1])
	if err != nil {
		hr = 0
	}
	today := time.Now().Day()
	day = day - today
	return day, hr
}

func (s *server) ListAccessPoints(in *pb.APRequest, stream pb.APService_ListAccessPointsServer) error {
	ts := in.Timestamp
	day, hr := getDayAndHourFromTimestamp(ts)

	apList := DBService.GetHistoryForAllAPs(day, hr)

	log.Printf("Sending %d APs ...", len(apList))

	for _, ap := range apList {
		location := locations[ap.Name]

		load, _ := strconv.Atoi(ap.Load)

		if err := stream.Send(
			&pb.APResponse{
				Accesspoint: &pb.AccessPoint{
					Name:      ap.Name,
					Lat:       location.lat,
					Long:      location.long,
					Intensity: int64(load),
					Max:       int64(ap.Max),
					Min:       int64(ap.Min),
				},
			}); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) ListAllAPNames(in *emptypb.Empty, stream pb.APService_ListAllAPNamesServer) error {
	names := DBService.GetAllNames()
	for _, name := range names {
		if err := stream.Send(
			&pb.APName{
				Name: name,
			}); err != nil {
			return err
		}
	}
	return nil
}

type location struct {
	lat  string
	long string
}

var locations map[string]location

func main() {

	// DBService.DeleteOldTables()
	// DBService.SetupHistoryTable()
	// DBService.SetupFutureTable()
	// DBService.UpdateHistory(0,0, 33, "apa08-1w4")


	// res := LRZscraper.GetHistoriesFrom(30)
	// LRZscraper.StoreHistories(res)

	// DBService.TestExample()
	// DBService.TestQuestionMark()
	// DBService.UpdateToday()

	// _ = LRZscraper.GetTodaysData()
	LRZscraper.Nothing()

	// DBService.UpdateTomorrow()

	// if true {
	// 	return
	// }

	locations = make(map[string]location)
	apList := DBService.RetrieveAPsOfTUM(true)
	for _, ap := range apList {
		locations[ap.Name] = location{ap.Lat, ap.Long}
	}

	startServer()

}

func startServer() {
	// host := "192.168.0.109"
	host := "172.20.10.7"
	port := 50053

	fmt.Println("Starting server...")

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAPServiceServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())

	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("%s:%d", host, port),
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

	grcpPort := 50054
	gwServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, grcpPort),
		Handler: gwmux,
	}

	log.Printf("Serving gRPC-Gateway on http://%s:%d", host, grcpPort)
	log.Fatalln(gwServer.ListenAndServe())
}

func forecast() {
	path := "./forecasting.py"
	cmd := exec.Command("python", "-u", path)
	out, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}
