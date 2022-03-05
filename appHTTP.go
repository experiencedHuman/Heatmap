package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	csvStr := "row1,value1,value2\nrow2,value1,value2"
	fmt.Println("csv content: \n", csvStr)
	req, err := http.NewRequest("GET", "http://graphite-kom.srv.lrz.de/render/?xFormat=%d.%m.%20%H:%M&tz=CET&from=-700days&target=cactiStyle(group(alias(ap.gesamt.ssid.eduroam,%22eduroam%22),alias(ap.gesamt.ssid.lrz,%22lrz%22),alias(ap.gesamt.ssid.mwn-events,%22mwn-events%22),alias(ap.gesamt.ssid.@BayernWLAN,%22@BayernWLAN%22),alias(ap.gesamt.ssid.other,%22other%22)))&format=csv", strings.NewReader(csvStr))
	if err != nil {
		fmt.Println("Error: could not get csv data from URL!")
	}
	results, _ := ReadCSVFromHttpRequest(req)
	fmt.Println("parsed csv: \n", results)
}

func ReadCSVFromHttpRequest(req *http.Request) ([][]string, error) {
	// parse POST body as csv
	reader := csv.NewReader(req.Body)
	var results [][]string
	for {
		// read one row from csv
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// add record to result set
		results = append(results, record)
	}
	return results, nil
}
