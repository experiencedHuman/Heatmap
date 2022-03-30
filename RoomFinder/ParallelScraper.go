package RoomFinder

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

type Coordinate struct {
	// visited   bool
	latitude  float64
	longitude float64
}

type Room struct {
	RoomNumber string
	BuildingNr string
}

var rooms map[string]Coordinate

func Scrape() {
	// url := "http://portal.mytum.de/displayRoomMap?"

	data := []Room{
		{RoomNumber: "5501", BuildingNr: "5509"},
		{RoomNumber: "5502", BuildingNr: "5508"},
	}

	// Instantiate default collector
	c := colly.NewCollector(
		colly.Async(true),
	)

	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	c.OnHTML("html", func(h *colly.HTMLElement) {
		roomNr := h.Request.URL
		buildingNr := h.Request.URL
		roomFinderID := fmt.Sprintf("%s@%s", roomNr, buildingNr)
		if _, ok := rooms[roomFinderID]; ok {
			// already visited
		} else {
			// not yet visited
			e := h.DOM.Find("a[href^='http://maps.google.com']")
			link, exists := e.Attr("href")
			if exists {
				lat, long := getLatLongFromURL(link)
				rooms[roomFinderID] = Coordinate{latitude: lat, longitude: long}
			}
		}
	})

	for _, val := range data {
		// Add URLs to the queue
		q.AddURL(fmt.Sprintf("http://portal.mytum.de/displayRoomMap?%s@%s", val.RoomNumber, val.BuildingNr))
	}
	// Consume URLs
	q.Run(c)
	c.Wait()
}

func getLatLongFromURL(url string) (float64, float64) {
	parts := strings.Split(url, "&")

	latLong := strings.Split(parts[0], "=")[1]
	latStr, longStr, _ := strings.Cut(latLong, ",")

	// spnLatLong := strings.Split(parts[1], "=")[1]
	// spnLatStr, spnLongStr, _ := strings.Cut(spnLatLong, ",")

	lat, _ := strconv.ParseFloat(latStr, 64)
	long, _ := strconv.ParseFloat(longStr, 64)

	// spnLat, _  := strconv.ParseFloat(spnLatStr, 64)
	// spnLong, _ := strconv.ParseFloat(spnLongStr, 64)

	return lat, long //, spnLat, spnLong
}

func scrapeBuildingNrFromAddress(address string) string {
	re := regexp.MustCompile("[0-9]+")
	buildingNr := re.FindString(address)
	return buildingNr
}
