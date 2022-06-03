package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/kvogli/Heatmap/proto/api"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kvogli/Heatmap/DBService"
	"github.com/kvogli/Heatmap/RoomFinder"
	"github.com/kvogli/Heatmap/LRZscraper"
)

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

// Retrieves all access points from the database
// and stores them in JSON format in `dst` e.g. "data/json/ap.json"
func saveAPsToJSON(dst string, totalLoad int) {
	APs := DBService.RetrieveAPsOfTUM(apstatDB, true)
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

	log.Printf("Received request for AP with name: %s", name)

	db := DBService.InitDB(ApstatTable)
	ap := DBService.RetrieveAccessPointByName(db, name)

	return &pb.AccessPoint{
		Name:      ap.Name,
		Lat:       ap.Lat,
		Long:      ap.Long,
		Intensity: ap.Load}, nil
}

func (s *server) ListAccessPoints(in *emptypb.Empty, stream pb.APService_ListAccessPointsServer) error {
	db := DBService.InitDB(ApstatTable)
	apList := DBService.RetrieveAPsOfTUM(db, true)

	log.Printf("Sending %d APs ...", len(apList))

	i := 1
	for _, ap := range apList {
		nty := fmt.Sprintf("%d", i)
		i++

		if err := stream.Send(
			&pb.APResponse{
				Accesspoint: &pb.AccessPoint{
					Name:      ap.Name,
					Lat:       ap.Lat,
					Long:      ap.Long,
					Intensity: nty},
			}); err != nil { //TODO implement intensity
			return err
		}
	}

	return nil
}

func main() {

	_ = "http://graphite-kom.srv.lrz.de/render/?width=640&height=240&title=SSIDs%20(weekly)&areaMode=stacked&xFormat=%25d.%25m.&tz=CET&from=-8days&target=cactiStyle(group(alias(ap.apa01-0lj.ssid.eduroam,%22eduroam%22),alias(ap.apa01-0lj.ssid.lrz,%22lrz%22),alias(ap.apa01-0lj.ssid.mwn-events,%22mwn-events%22),alias(ap.apa01-0lj.ssid.@BayernWLAN,%22@BayernWLAN%22),alias(ap.apa01-0lj.ssid.other,%22other%22)))&fontName=Courier&format=csv"
	// LRZscraper.GetGraphiteData("data/csv/apa01-0lj.csv", url)
	// LRZscraper.GetGraphiteDataAsJSON("apa01-1mo", "")
	res := LRZscraper.GetGraphiteDataAsJSON("apa02-1mo", "")

	for _, val := range res {
		fmt.Println(val.Datapoints)
	}

	if true {
		return
	}

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
