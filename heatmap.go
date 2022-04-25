package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"encoding/csv"
	"encoding/json"
	"net/http"

	"github.com/kvogli/Heatmap/DBService"
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
	Lat       string
	Long      string
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
		jsonEntry := JsonEntry{intensity, ap.Lat, ap.Long, ap.Floor}
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

func main() {
	db := DBService.InitDB(ApstatTable)
	// DBService.AddColumn(db, "apstat", "Floor")
	// DBService.UpdateColumn("apstat", "Long", "longitude", " IS NULL")

	// DBService.UpdateColumn("apstat", "Lat", "lat", "Lat != 'NULL'")
	// DBService.UpdateColumn("apstat", "Long", "long", "Long != 'zrf'")
	// scrapeRoomFinder()
	APs1 := DBService.RetrieveAPs(db, false)
	APs2 := DBService.RetrieveAPs(db, true)
	// add floors
	for _, ap := range APs1 {
		floor := string(ap.Name[6])
		where := fmt.Sprintf("ID='%s'", ap.ID)
		DBService.UpdateColumn(db, "apstat", "Floor", floor, where)
	}
	for _, ap := range APs2 {
		floor := string(ap.Name[6])
		where := fmt.Sprintf("ID='%s'", ap.ID)
		DBService.UpdateColumn(db, "apstat", "Floor", floor, where)
	}
}

func scrapeRoomFinder() int {
	db := DBService.InitDB(ApstatTable)
	APs := DBService.RetrieveAPs(db, false)
	roomInfos, totalLoad := RoomFinder.PrepareDataToScrape(APs)
	res := RoomFinder.ScrapeURLs(roomInfos)

	for key, val := range res {
		where := fmt.Sprintf("ID='%s'", key)
		DBService.UpdateColumn(db, "apstat", "Lat", val.Lat, where)
		DBService.UpdateColumn(db, "apstat", "Long", val.Long, where)
	}

	return totalLoad
}
