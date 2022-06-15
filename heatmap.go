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

	// DBService.SetupHistoryTable()

	// 1
	// LRZscraper.StoreHistory()

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
