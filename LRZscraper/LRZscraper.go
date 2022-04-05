package LRZscraper

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

// It scrapes the table data from URL https://wlan.lrz.de/apstat/ and
// stores the scraped data in csv format under the path parameter  'fName'
func ScrapeApstat(fName string) {
	apstatURL := "https://wlan.lrz.de/apstat/"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Error: could not create file! %q", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	c := colly.NewCollector(
		colly.AllowedDomains("wlan.lrz.de"),
		colly.Debugger(&debug.LogDebugger{}),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// it uses jQuery selectors to scrape the table with id "aptable" row by row
	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.DOM.Find("table#aptable > tbody > tr").Each(func(i int, s *goquery.Selection) {
			address := s.ChildrenFiltered("td:nth-child(1)").Text()
			room := s.ChildrenFiltered("td:nth-child(2)").Text()
			apName := s.ChildrenFiltered("td:nth-child(3)").Text()
			apStatus := s.ChildrenFiltered("td:nth-child(4)").Text()
			apStatus = strings.TrimSpace(apStatus)
			apType := s.ChildrenFiltered("td:nth-child(5)").Text()
			load := s.ChildrenFiltered("td:nth-child(6)").Text()
			writer.Write([]string{
				address,
				room,
				apName,
				apStatus,
				apType,
				load,
			})
		})
	})

	c.Visit(apstatURL)
}

// It scrapes the table data from URL https://wlan.lrz.de/apstat/ublist/ and
// stores the scraped data in csv format under the path parameter 'fName'
func ScrapeApstatUblist(fName string) {
	ublistURL := "https://wlan.lrz.de/apstat/ublist/"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Error: could not create file! %q", err)
		return
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

	c.OnError(func(request *colly.Response, err error) {
		log.Println("Request URL:", request.Request.URL, "failed with response:", request, "\nError:", err)
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
