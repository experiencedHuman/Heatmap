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

var rooms map[string]Coordinate

// This function generates a number of URLs (based on roomIDs) and 
// visits them all in parallel and scrapes the coordinates for each room.
// It returns a map of all rooms with key being roomID, and value room coordinates
// Reference -> If you want to refer to a room map from outside the portal, please use the complete path:
// http://portal.mytum.de/displayRoomMap?roomnumber@builingnumber
func Scrape(roomIDs []string) {
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

	rooms := make(map[string]Coordinate)

	c.OnHTML("html", func(h *colly.HTMLElement) {
		roomRE := regexp.MustCompile("[0-9]+")
		roomNr := roomRE.FindString(h.Request.URL.String())

		buildingRE := regexp.MustCompile("[0-9]+$")
		buildingNr := buildingRE.FindString(h.Request.URL.String())
		
		roomFinderID := fmt.Sprintf("%s@%s", roomNr, buildingNr)
		if _, ok := rooms[roomFinderID]; ok {
			// already visited
			// fmt.Println("if", roomFinderID)
		} else {
			// not yet visited
			e := h.DOM.Find("a[href^='http://maps.google.com']")
			link, exists := e.Attr("href")
			if exists {
				lat, long := getLatLongFromURL(link)
				// fmt.Println("else", roomFinderID)
				// fmt.Println(rooms[roomFinderID])
				rooms[roomFinderID] = Coordinate{latitude: lat, longitude: long}
			}
		}
	})

	for _, roomID := range roomIDs {
		// Add URLs to the queue
		q.AddURL(fmt.Sprintf("http://portal.mytum.de/displayRoomMap?%s", roomID))
	}
	// Consume URLs
	q.Run(c)
	c.Wait()
	fmt.Println("finito")
}

// This function scrapes latitude and longitude from url and 
// returns them as float64 values
func getLatLongFromURL(url string) (float64, float64) {
	parts := strings.Split(url, "&")

	latLong := strings.Split(parts[0], "=")[1]
	latStr, longStr, _ := strings.Cut(latLong, ",")

	lat, _ := strconv.ParseFloat(latStr, 64)
	long, _ := strconv.ParseFloat(longStr, 64)

	return lat, long
}

// This function scrapes the buildingnumber from a long address description and returns it
// Reference -> If you want to refer to a room map from outside the portal, please use the complete path:
// http://portal.mytum.de/displayRoomMap?roomnumber@builingnumber
func scrapeBuildingNrFromAddress(address string) string {
	re := regexp.MustCompile("[0-9]+")
	buildingNr := re.FindString(address)
	return buildingNr
}

// This function scrapes the roomnumber from a longer room description and returns it
// Reference -> If you want to refer to a room map from outside the portal, please use the complete path:
// http://portal.mytum.de/displayRoomMap?roomnumber@builingnumber
func scrapeRoomNrFromRoomName(roomName string) string {
	re := regexp.MustCompile("[0-9]+.[0-9]+(.[0-9])?")
	roomNr := re.FindString(roomName)
	return roomNr
}

// This function prepares the roomIDs that will later be concatenated to the RoomFinder's url/path
// to get the room map and then scrape its coordinate. It returns a slice of all roomIDs
func PrepareDataToScrape() []string {
	db := InitDB("./overview.db")
	fmt.Println("Preparing data")
	res := ReadItem(db)
	var data []string
	for _, val := range res {
		roomNr := scrapeRoomNrFromRoomName(val.Room)
		buildingNr := scrapeBuildingNrFromAddress(val.Address)
		roomID := fmt.Sprintf("%s@%s", roomNr, buildingNr)
		data = append(data, roomID)
	}
	return data
}
