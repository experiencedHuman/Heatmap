package main

import (
	// "database/sql"
	_ "github.com/mattn/go-sqlite3"

	"io"
	"os"
	"fmt"
	"log"
	"strings"
	"net/http"
	"encoding/csv"

	"github.com/experiencedHuman/heatmap/LRZscraper"
)

func readCSVFromUrl(url string) {
	resp, httpError := http.Get(url)
	if httpError != nil {
		fmt.Println("Could not retrieve csv data from URL!", httpError)
		return
	} else {
		fmt.Println("Successfully retrieved data from URL!")
	}

	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)
	csvReader.Comma = ','

	fName := "graphData.csv"
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
			name 	:= fields[0]
			current := strings.Split(fields[1], ":")[1]
			max 	:= strings.Split(fields[2], ":")[1]
			min 	:= strings.Split(fields[3], ":")[1]

			dateAndTime := data[1]
			other := data[2] // TODO find out what other is
			csvWriter.Write([]string{
				name, current, max, min, dateAndTime, other,
			})
		}
	}
}

func main() {
	LRZscraper.ScrapeListOfSubdistricts()
	LRZscraper.ScrapeOverviewOfAPs()
	
	// url := "http://graphite-kom.srv.lrz.de/render/?xFormat=%d.%m.%20%H:%M&tz=CET&from=-2days&target=cactiStyle(group(alias(ap.gesamt.ssid.eduroam,%22eduroam%22),alias(ap.gesamt.ssid.lrz,%22lrz%22),alias(ap.gesamt.ssid.mwn-events,%22mwn-events%22),alias(ap.gesamt.ssid.@BayernWLAN,%22@BayernWLAN%22),alias(ap.gesamt.ssid.other,%22other%22)))&format=csv"
	// readCSVFromUrl(url)

	// // open a local database instance, located in this directory
	// db, _ := sql.Open("sqlite3", "./accesspoints.db")

	// stmt, _ := db.Prepare(`
	// 	INSERT INTO accesspoints (address) values (?)
	// `)

	// // stmt, _ := db.Prepare(`
	// // 	CREATE TABLE IF NOT EXISTS "accesspoints" (
	// // 		"ID"	INTEGER NOT NULL PRIMARY KEY,
	// // 		"address"	TEXT
	// // 	);
	// // `)

	// stmt.Exec("Garching")
}
