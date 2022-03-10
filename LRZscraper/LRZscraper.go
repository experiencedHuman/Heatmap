package LRZscraper

import (
	"encoding/csv"
    "fmt"
    "log"
    "os"
	"github.com/gocolly/colly"
)

func ScrapeData() {
	fName := "data.csv"
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
	
    c.OnHTML("table", func(e *colly.HTMLElement) {
        e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
            writer.Write([]string{
                el.ChildText("td:nth-child(1)"),
                el.ChildText("td:nth-child(2)"),
                el.ChildText("td:nth-child(3)"),
				el.ChildText("td:nth-child(4)"),
            })
        })
        fmt.Println("Scrapping Complete")
    })

	c.Visit("https://wlan.lrz.de/apstat/ublist/")
}
