package LRZscraper

import (
	"fmt"
	"encoding/csv"
	"log"
	"os"
	"github.com/gocolly/colly"
)

// It scrapes the table data of a single subdistrict's (Unterbezirk auf Deutsch) page
// from "https://wlan.lrz.de/apstat/ublist/"
func ScrapeListOfSubdistricts(fName string) {
	subdistrictsURL := "https://wlan.lrz.de/apstat/ublist/"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Could not create file, err: %q", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	c := colly.NewCollector(
		colly.AllowedDomains("wlan.lrz.de"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(request *colly.Response, err error) {
	    fmt.Println("Request URL:", request.Request.URL, "failed with response:", request, "\nError:", err)
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

	c.Limit(&colly.LimitRule{DomainGlob: "wlan.lrz.de", Parallelism: 2})
	error := c.Visit(subdistrictsURL)
	if error != nil {
		fmt.Println(error)
	}
	c.Wait()
}

func ScrapeOverviewOfAPs(filename string) {
	overviewURL := "https://wlan.lrz.de/apstat"
	fName := filename
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Could not create file, err: %q", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	c := colly.NewCollector(
		colly.AllowedDomains("wlan.lrz.de"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnHTML("tbody > tr", func(e *colly.HTMLElement) {
		selection := e.DOM.Find("td:first-child > a[href]").Text()
		// ss := fmt.Sprintf("%d", selection)
		writer.Write([]string{
			selection,
		})
	})

	c.Visit(overviewURL)
}
