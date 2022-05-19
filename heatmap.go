package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/experiencedHuman/heatmap/LRZscraper"
	"google.golang.org/grpc"

	pb "github.com/experiencedHuman/heatmap/api"
)

func readCSVFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Could not retrieve csv data from URL!")
		return nil, err
	} else {
		fmt.Println("Successfully retrieved data from URL!")
	}

	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)
	csvReader.Comma = ','
	var data []string
	for i := 0; i < 5; i++ {
		data, err = csvReader.Read()
		if err != nil {
			return nil, err
		} else {
			fmt.Printf("%+v\n", data)
		}
	}
	return make([][]string, 0), nil
}

func getGraphiteDataFromLRZ() {
	url := "http://graphite-kom.srv.lrz.de/render/?xFormat=%d.%m.%20%H:%M&tz=CET&from=-5days&target=cactiStyle(group(alias(ap.gesamt.ssid.eduroam,%22eduroam%22),alias(ap.gesamt.ssid.lrz,%22lrz%22),alias(ap.gesamt.ssid.mwn-events,%22mwn-events%22),alias(ap.gesamt.ssid.@BayernWLAN,%22@BayernWLAN%22),alias(ap.gesamt.ssid.other,%22other%22)))&format=csv"

	data, err := readCSVFromUrl(url)
	if err != nil {
		panic(err)
	}

	for idx, row := range data {
		// skip header
		if idx == 0 {
			continue
		}

		if idx == 6 {
			break
		}

		fmt.Println(row[2])
	}
}

type APServiceServer struct {
	pb.UnimplementedAPServiceServer
}

func (s *APServiceServer) GetAccessPoint(ctx context.Context, in *pb.Empty) (*pb.AccessPoint, error) {
	log.Println("Received request from client! \n Returning \"apa99-k99\" as a sample response!")
	return &pb.AccessPoint{Name: "apa99-k99"}, nil
}

func (s *APServiceServer) ListAccessPoints(in *pb.Empty, stream pb.APService_ListAccessPointsServer) error {
	// TODO: get actual data; merge repo from master branch
	for i := 0; i < 10; i++ {
		n := fmt.Sprintf("%d", i)
		if err:= stream.Send(&pb.AccessPoint{Name: "apa", Lat: n, Long: n, Intensity: n}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if false {
		LRZscraper.ScrapeData()
	}
	
	port := 50051
	
	fmt.Println("Starting server...")
	
	lis, err := net.Listen("tcp", fmt.Sprintf("192.168.0.109:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAPServiceServer(s, &APServiceServer{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
