package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/experiencedHuman/heatmap/LRZscraper"
)

func getDataFromURL(fName, url string) {
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
			fields 	:= strings.Fields(data[0]) // get substrings separated by whitespaces
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

// for csv/graphData.csv
func storeDataInSQLite(dbPath string) {
	csvFile, err := os.Open("csv/graphData.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'
	data, err := csvReader.ReadAll() // TODO use csvReader.Read() for big files
	if err != nil {
		log.Fatal(err)
	}

	// open a local database instance, located in dbPath
	db, sqlError := sql.Open("sqlite3", dbPath)
	if sqlError != nil {
		panic(sqlError)
	}

	stmt, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "accesspoints" (
			"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"network"	TEXT,
			"current"	TEXT,
			"max"		TEXT,
			"min"		TEXT,
			"other" 	TEXT
		);
	`)
	stmt.Exec()

	// fmt.Println(network, current, max, min, other)
	stmt, dbError := db.Prepare(`
		INSERT INTO accesspoints (network, current, max, min, other) values (?,?,?,?,?)
	`)

	if dbError != nil {
		panic(dbError)
	}

	fmt.Println("Storing data in SQLite ...")
	for r := range data {
		network := data[r][0]
		current := data[r][1]
		max 	:= data[r][2]
		min 	:= data[r][3]
		other 	:= data[r][4]

		_, err := stmt.Exec(network, current, max, min, other)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Finished data storing")
}

// for csv/overview.csv
func storeOverviewInSQLite(dbPath string) {
	csvFile, err := os.Open("csv/overview.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'
	data, err := csvReader.ReadAll() // TODO use csvReader.Read() for big files
	if err != nil {
		log.Fatal(err)
	}

	// open a local database instance, located in dbPath
	db, sqlError := sql.Open("sqlite3", dbPath)
	if sqlError != nil {
		panic(sqlError)
	}

	stmt, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "overview" (
			"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"address"	TEXT,
			"roomNr"	TEXT,
			"apName"	TEXT,
			"status"	TEXT,
			"apType"	TEXT,
			"apLoad" 	TEXT
		);
	`)
	stmt.Exec()

	stmt, dbError := db.Prepare(`
		INSERT INTO overview (address, roomNr, apName, status, apType, apLoad) values (?,?,?,?,?,?)
	`)

	if dbError != nil {
		panic(dbError)
	}

	fmt.Println("Storing data in SQLite ...")
	for r := range data {
		address := data[r][0]
		roomNr 	:= data[r][1]
		apName 	:= data[r][2]
		status 	:= data[r][3]
		apType 	:= data[r][4]
		apLoad 	:= data[r][5]

		_, err := stmt.Exec(address, roomNr, apName, status, apType, apLoad)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Finished data storing")
}

// func getListOfIntensities() []int {

// }

// func saveApLoadToJsonFile() {

// }

func main() {
	// LRZscraper.ScrapeListOfSubdistricts("csv/subdistricts.csv")
	// LRZscraper.ScrapeOverviewOfAPs("csv/overview.csv")
	// storeOverviewInSQLite("./overview.db")

	// url := "http://graphite-kom.srv.lrz.de/render/?xFormat=%d.%m.%20%H:%M&tz=CET&from=-2days&target=cactiStyle(group(alias(ap.gesamt.ssid.eduroam,%22eduroam%22),alias(ap.gesamt.ssid.lrz,%22lrz%22),alias(ap.gesamt.ssid.mwn-events,%22mwn-events%22),alias(ap.gesamt.ssid.@BayernWLAN,%22@BayernWLAN%22),alias(ap.gesamt.ssid.other,%22other%22)))&format=csv"
	// getDataFromURL("csv/graphData.csv", url)

	// storeDataInSQLite("./accesspoints.db")
	LRZscraper.ScrapeMapCoordinatesForRoom("1", "5406")

}
