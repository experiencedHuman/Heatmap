package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kvogli/Heatmap/LRZscraper"
	"github.com/kvogli/Heatmap/RoomFinder"

	"os"
	"strings"
)

func getDataFromURL(fName, url string) {
	resp, httpError := http.Get(url)
	if httpError != nil {
		log.Println("Could not retrieve csv data from URL!", httpError)
		return
	} else {
		log.Println("Successfully retrieved data from URL!")
	}

	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)
	csvReader.Comma = ','

	file, osError := os.Create(fName)
	if osError != nil {
		log.Fatalf("Could not create file, err: %q", osError)
		return
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Comma = ';'
	defer csvWriter.Flush()

	var data []string
	for i := 0; ; i++ {
		data, httpError = csvReader.Read()
		if httpError == io.EOF {
			break
		} else if httpError != nil {
			panic(httpError)
		} else {
			fields := strings.Fields(data[0]) // get substrings separated by whitespaces
			name := fields[0]
			current := strings.Split(fields[1], ":")[1]
			max := strings.Split(fields[2], ":")[1]
			min := strings.Split(fields[3], ":")[1]

			dateAndTime := data[1]
			other := data[2] // TODO find out what other is
			csvWriter.Write([]string{
				name, current, max, min, dateAndTime, other,
			})
		}
	}
}

type AccessPoint struct {
	Intensity float64
	Latitude  float64
	Longitude float64
}

func saveApLoadToJsonFile() {
	roomInfos, totalLoad := RoomFinder.PrepareDataToScrape()
	coordinatesMap := RoomFinder.Scrape(roomInfos) // TODO use ignore result, which is a map of room ids and coordinates

	accessPoints := make([]AccessPoint, 0)
	for _, roomInfo := range roomInfos {
		var intensity float64 = float64(roomInfo.RoomLoad) / float64(totalLoad)
		if coord, exists := coordinatesMap[roomInfo.RoomID]; exists {
			ap := AccessPoint{Intensity: intensity, Latitude: coord.Lat, Longitude: coord.Long}
			accessPoints = append(accessPoints, ap)
		}
	}

	bytes, err := json.Marshal(accessPoints)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("accessPoints.json", bytes, 0644)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Data was saved successfully to accessPoints.json")
	}
}

func main() {
	// LRZscraper.ScrapeApstat("data/csv/apstat.csv")
	// LRZscraper.StoreApstatInSQLite("data/sqlite/apstat.db")
	// saveApLoadToJsonFile()

	// res := LRZscraper.FetchApstatData("apstat")
	// println(len(res))
	 
	// LRZscraper.ScrapeApstat("data/csv/apstat.csv")
	// LRZscraper.StoreApstatInSQLite("apstat")
	LRZscraper.PopulateNewColumn("apstat", "RF_ID")
}
