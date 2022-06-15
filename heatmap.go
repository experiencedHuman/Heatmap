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

	// "time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/kvogli/Heatmap/proto/api"

	// "google.golang.org/protobuf/types/known/emptypb"

	"github.com/kvogli/Heatmap/DBService"
	// "github.com/kvogli/Heatmap/LRZscraper"
	"github.com/kvogli/Heatmap/RoomFinder"
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

type server struct {
	pb.UnimplementedAPServiceServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) GetAccessPoint(ctx context.Context, in *pb.APRequest) (*pb.AccessPoint, error) {
	name := in.Name
	ts := in.Timestamp
	log.Printf("Received request for AP with name: %s and timestamp: %s", name, in.Timestamp)

	db := DBService.InitDB(heatmapDB)
	// ap := DBService.RetrieveAccessPointByName(db, name)
	day, hr := getDayAndHourFromTimestamp(ts)
	ap := DBService.GetHistoryOfSingleAP(name, day, hr)
	db.Close()

	return &pb.AccessPoint{
		Name:      ap.Name,
		Lat:       ap.Lat,
		Long:      ap.Long,
		Intensity: ap.Load}, nil
}

func getDayAndHourFromTimestamp(timestamp string) (int, int) {
	ts := strings.Split(timestamp, " ")
	date := ts[0]
	ymd := strings.Split(date, "-")
	day, err := strconv.Atoi(ymd[2])
	if err != nil {
		day = 0
	}
	hr, err := strconv.Atoi(ts[1])
	if err != nil {
		hr = 0
	}
	return day, hr
}

func (s *server) ListAccessPoints(in *pb.APRequest, stream pb.APService_ListAccessPointsServer) error {
	ts := in.Timestamp
	day, hr := getDayAndHourFromTimestamp(ts)
	apList := DBService.GetHistoryOfAllAccessPoints(day, hr)

	log.Printf("Sending %d APs ...", len(apList))

	i := 1
	for _, ap := range apList {
		nty := fmt.Sprintf("%d", i)
		i++

		if err := stream.Send(
			&pb.APResponse{
				Accesspoint: &pb.AccessPoint{
					Name:      ap.Name,
					Lat:       ap.Lat,  // TODO handle nil value
					Long:      ap.Long, // TODO handle nil value
					Intensity: nty},
			}); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// TODO migrate all tables from different DB into heatmapDB. [done]

	// TODO restructure DB [done]
	// TODO code refactor [db done]

	// TODO handle Different lengths:
	// TODO handle Skipping ... Bad JSON response form server
	
	// TODO store data in DB
	// TODO store data for the day mapping, store max and min avg values of each ap
	
	// TODO query history db from client

	// TODO call python script and check its output (csv file) from golang [half done]
	// TODO STORE forecasting csv in forecasting table in the DB
	
	// TODO improve forecasting

	// TODO update proto & implement new functionality

	//Could not decode JSON response: invalid character 'T' looking for beginning of value
	//http://graphite-kom.srv.lrz.de/render/?width=640&height=240&title=&title=SSIDs%20(weekly)&areaMode=stacked&xFormat=%25d.%25m.&tz=CET&from=-30days&target=cactiStyle(group(alias(ap.apa05-1mm.ssid.eduroam,%22eduroam%22),alias(ap.apa05-1mm.ssid.lrz,%22lrz%22),alias(ap.apa05-1mm.ssid.mwn-events,%22mwn-events%22),alias(ap.apa05-1mm.ssid.@BayernWLAN,%22@BayernWLAN%22),alias(ap.apa05-1mm.ssid.other,%22other%22)))&fontName=Courier&format=json
	//http://graphite-kom.srv.lrz.de/render/?width=640&height=240&title=&title=SSIDs%20(weekly)&areaMode=stacked&xFormat=%25d.%25m.&tz=CET&from=-30days&target=cactiStyle(group(alias(ap.apa04-1w6.ssid.eduroam,%22eduroam%22),alias(ap.apa04-1w6.ssid.lrz,%22lrz%22),alias(ap.apa04-1w6.ssid.mwn-events,%22mwn-events%22),alias(ap.apa04-1w6.ssid.@BayernWLAN,%22@BayernWLAN%22),alias(ap.apa04-1w6.ssid.other,%22other%22)))&fontName=Courier&format=json

	//apa03-2qu
	// apa09-1w6 -> different lengths 8640 vs 2880
	// LRZscraper.SaveLast30DaysForAP(
	// 	DBService.AccessPoint{Name: "apa05-0mg"})

	// ap := DBService.GetApDataFromLast30Days("apa05-0mg", 5, 12)
	// fmt.Printf("Network load = %s", ap.Load)

	// -------------------------
	// list := DBService.GetHistoryOfAllAccessPoints(5, 20)
	// for _, ap := range list {
	// 	fmt.Println(ap.Name, ap.Load)
	// }
	// fmt.Println(len(list))
	// -------------------------

	// ap := DBService.AccessPoint{Name: "apa05-0mg"}

	// DBService.SetupHistoryTable()

	// 1
	// LRZscraper.StoreHistory()

	// DBService.HistoryCSVtoSQLite()

	// DBService.AddColumn("history", "Max", "INTEGER")
	// DBService.AddColumn("history", "Min", "INTEGER")
	DBService.JoinMaxMin()

	// 2
	// DBService.AddColumn("apstat", "Max", "INTEGER")
	// DBService.AddColumn("apstat", "Min", "INTEGER")

	// forecast()

}

func startServer() {
	host := "192.168.0.109"
	port := 50051

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

	gwServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, 50052),
		Handler: gwmux,
	}

	log.Printf("Serving gRPC-Gateway on http://%s:%d", host, 50052)
	log.Fatalln(gwServer.ListenAndServe())
}

func forecast() {
	path := "./forecast-script.py"
	cmd := exec.Command("python", "-u", path)
	out, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}
