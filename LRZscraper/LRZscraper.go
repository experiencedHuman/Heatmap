package LRZscraper

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

type AP struct {
	address string
	room    string
	name    string
	status  string
	typ     string
	load    string
}

func Nothing() {
	
}

// It scrapes the html table data from "https://wlan.lrz.de/apstat/""
func ScrapeApstat(filename string) []AP {
	apstatURL := "https://wlan.lrz.de/apstat/"

	c := colly.NewCollector(
		colly.AllowedDomains("wlan.lrz.de"),
		colly.Debugger(&debug.LogDebugger{}),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	aps := make([]AP, 0)

	// it uses jQuery selectors to scrape the table with id "aptable" row by row
	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.DOM.Find("table#aptable > tbody > tr").Each(func(i int, s *goquery.Selection) {
			// scrape overview data for each access point
			address := s.ChildrenFiltered("td:nth-child(1)").Text()
			room := s.ChildrenFiltered("td:nth-child(2)").Text()
			apName := s.ChildrenFiltered("td:nth-child(3)").Text()
			apStatus := s.ChildrenFiltered("td:nth-child(4)").Text()
			apStatus = strings.TrimSpace(apStatus)
			apType := s.ChildrenFiltered("td:nth-child(5)").Text()
			load := s.ChildrenFiltered("td:nth-child(6)").Text()
			// store scraped data to csv file
			ap := AP{address, room, apName, apStatus, apType, load}
			aps = append(aps, ap)
		})
	})

	c.Visit(apstatURL)
	c.Wait()

	return aps
}

func getTotAvg(load string) string {
	rgx := regexp.MustCompile(`\((.*?)\)`)
	loadStr := rgx.FindStringSubmatch(load)
	loads := strings.Split(loadStr[1], " - ")
	totAvg := loads[3]
	return totAvg
}


// Stores the scraped data in csv format under the destination path 'filename'.
func StoreApstatInCSV(aps []AP, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("System error: could not create file! %q", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	for _, ap := range aps {
		writer.Write([]string{
			ap.address,
			ap.room,
			ap.name,
			ap.status,
			ap.typ,
			ap.load,
		})
	}
}

// It scrapes the html table data from "https://wlan.lrz.de/apstat/ublist/"" and
// stores the scraped data in csv format under the path parameter 'filename'
func ScrapeApstatUblist(filename string) {
	ublistURL := "https://wlan.lrz.de/apstat/ublist/"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("System error: could not create file! %q", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	c := colly.NewCollector(
		colly.AllowedDomains("wlan.lrz.de"),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnError(func(res *colly.Response, err error) {
		log.Println("Request URL:", res.Request.URL, "failed with response:", res, "\nError:", err)
	})

	// scrape table's head
	c.OnHTML("thead", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText("th:nth-child(1)"),
			"Link",
			e.ChildText("th:nth-child(2)"),
			e.ChildText("th:nth-child(3)"),
			e.ChildText("th:nth-child(4)"),
		})
	})

	// scrape table's body
	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			writer.Write([]string{
				el.ChildText("td:nth-child(1)"),
				el.ChildAttr("a", "href"), // the link of the address
				el.ChildText("td:nth-child(2)"),
				el.ChildText("td:nth-child(3)"),
				el.ChildText("td:nth-child(4)"),
			})
		})
	})

	c.Visit(ublistURL)
}
