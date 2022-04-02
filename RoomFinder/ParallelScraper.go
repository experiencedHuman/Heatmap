package RoomFinder

import (
	"log"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	// "github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/queue"
	"github.com/experiencedHuman/heatmap/LRZscraper"
)

type muxCache struct {
	mux 	sync.Mutex
	rooms map[string]Coordinate
}

var roomCache = muxCache{rooms: make(map[string]Coordinate)}

type Coordinate struct {
	// visited   bool
	Latitude  float64
	Longitude float64
}

func (cache *muxCache) setVisited(h *colly.HTMLElement, roomFinderID string) bool {
	cache.mux.Lock()
	defer cache.mux.Unlock()
	if _, ok := roomCache.rooms[roomFinderID]; ok {
		// already visited
		return true
	} else {
		// not yet visited
		e := h.DOM.Find("a[href^='http://maps.google.com']")
		link, exists := e.Attr("href")
		if exists {
			lat, long := getLatLongFromURL(link)
			roomCache.rooms[roomFinderID] = Coordinate{Latitude: lat, Longitude: long}
		}
		return false
	}
}

func showStatus(q *queue.Queue) {
	qSize, err := q.Size()
	if err != nil {
		log.Println("Error reading queue size!", err)
	} else {
			log.Println("Queue size:", qSize)
	}
}

// This function generates a number of URLs (based on roomIDs) and
// visits them all in parallel and scrapes the coordinates for each room.
// It returns a map of all rooms with key being roomID, and value room coordinates
// Reference -> If you want to refer to a room map from outside the portal, please use the complete path:
// http://portal.mytum.de/displayRoomMap?roomnumber@builingnumber
func Scrape(roomInfos []RoomInfo) map[string]Coordinate {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.Async(true),
		// colly.Debugger(&debug.LogDebugger{}),
	)

	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 20})

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		log.Println("Finished scraping", r.Request.URL)
	})


	c.OnHTML("html", func(h *colly.HTMLElement) {
		roomRE := regexp.MustCompile("[0-9]+")
		roomNr := roomRE.FindString(h.Request.URL.String())

		buildingRE := regexp.MustCompile("[0-9]+$")
		buildingNr := buildingRE.FindString(h.Request.URL.String())

		roomFinderID := fmt.Sprintf("%s@%s", roomNr, buildingNr)
		roomCache.setVisited(h, roomFinderID)
	})

	for _, roomInfo := range roomInfos {
		// Add URLs to the queue
		q.AddURL(fmt.Sprintf("http://portal.mytum.de/displayRoomMap?%s", roomInfo.RoomID))
	}
	start := time.Now()
	log.Println("Time now", start)
	// Consume URLs
	q.Run(c)
	c.Wait()
	elapsed := time.Since(start)
	// TODO add progress indicator
	// go showStatus(q)
	fmt.Println("Finished scraping locations. Time elapsed:", elapsed)
	return roomCache.rooms
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
	// TODO measure success rate of this regex
	roomNr := re.FindString(roomName)
	return roomNr
}

type RoomInfo struct {
	RoomID   string
	RoomLoad int
}

// This function prepares the roomIDs that will later be concatenated to the RoomFinder's url/path
// to get the room map and then scrape its coordinate. It returns a slice of all roomIDs
func PrepareDataToScrape() ([]RoomInfo, int) {
	db := LRZscraper.InitDB("./overview.db")
	fmt.Println("Preparing data")
	res := LRZscraper.ReadItem(db)
	var data []RoomInfo
	var total int
	for _, val := range res {
		roomNr := scrapeRoomNrFromRoomName(val.Room)
		buildingNr := scrapeBuildingNrFromAddress(val.Address)
		currTotalLoad := getCurrentTotalLoad(val.Load)
		total += currTotalLoad
		roomID := fmt.Sprintf("%s@%s", roomNr, buildingNr)
		data = append(data, RoomInfo{RoomID: roomID, RoomLoad: currTotalLoad})
	}
	return data, total
}

func getCurrentTotalLoad(load string) int {
	// this regex must match a substring beginning with '(', ignores first number and '-', and then gets the second number
	// based on: https://stackoverflow.com/questions/46817073/regular-expression-for-getting-numbers-in-between
	re := regexp.MustCompile(`\(\s*\d+\s*-\s*(\d+)`)
	match := re.FindStringSubmatch(load)
	currentLoad, err := strconv.Atoi(match[1]) // TODO error handling
	if err != nil {
		return 0 // TODO handle edge cases
	} else {
		return currentLoad
	}
}
