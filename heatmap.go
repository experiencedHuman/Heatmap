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
	"github.com/kvogli/Heatmap/NavigaTUM"
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

	result, _ := scrapeRoomFinder() //Note that room finder must first be scraped to jump to navigatum this way
	cnt := scrapeNavigaTUM(result)
	fmt.Println(cnt)
	
	// db := DBService.InitDB(ApstatTable)
	
	// APs1 := DBService.RetrieveAPs(db, false)
	// APs2 := DBService.RetrieveAPs(db, true)
	// // add floors
	// for _, ap := range APs1 {
	// 	floor := string(ap.Name[6])
	// 	where := fmt.Sprintf("ID='%s'", ap.ID)
	// 	DBService.UpdateColumn(db, "apstat", "Floor", floor, where)
	// }
	// for _, ap := range APs2 {
	// 	floor := string(ap.Name[6])
	// 	where := fmt.Sprintf("ID='%s'", ap.ID)
	// 	DBService.UpdateColumn(db, "apstat", "Floor", floor, where)
	// }

	// lat, long, _ := NavigaTUM.GetRoomCoordinates("9377")
	// fmt.Println(lat, long)
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
