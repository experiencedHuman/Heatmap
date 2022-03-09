package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"log"
	"os"
	"github.com/gocolly/colly"
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

func main() {
	// getGraphiteDataFromLRZ()

	fName := "data.csv"
	file, err := os.Create(fName)

	if err != nil {
		log.Fatalf("Could not create file, err: %q", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	c := colly.NewCollector()
	c.OnHTML("table#customers", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			writer.Write([]string{
				el.ChildText("td:nth-child(1)"),
				el.ChildText("td:nth-child(2)"),
				el.ChildText("td:nth-child(3)"),
			})
		})
		fmt.Println("Scrapping Complete")
	})
	c.Visit("https://www.w3schools.com/html/html_tables.asp")

}
