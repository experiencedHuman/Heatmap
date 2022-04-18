package RoomFinder

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/kvogli/Heatmap/DBService"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type muxCache struct {
	mux   sync.Mutex
	rooms map[string]Coordinate
}

var roomCache = muxCache{rooms: make(map[string]Coordinate)}

type Coordinate struct {
	exact bool // whether the coordinates are that of the room, or of the building
	Lat   float64
	Long  float64
}

func (cache *muxCache) setVisited(h *colly.HTMLElement, roomFinderID string) bool {
	cache.mux.Lock()
	defer cache.mux.Unlock()
	if _, ok := roomCache.rooms[roomFinderID]; ok {
		// already visited
		// TODO adjust intensity of the room/ap
		return true
	} else {
		// not yet visited
		e := h.DOM.Find("a[href^='http://maps.google.com']")
		link, exists := e.Attr("href")
		if exists {
			lat, long := getLatLongFromURL(link)
			if h.DOM.Find(".message").Length() == 0 {
				roomCache.rooms[roomFinderID] = Coordinate{exact: true, Lat: lat, Long: long}
			} else {
				/* Die angezeigte Position zeigt das Gebäude.
				Die Position des Raums innerhalb des Gebäudes ist leider nicht bekannt. */
				roomCache.rooms[roomFinderID] = Coordinate{exact: false, Lat: lat, Long: long}
			}
		}
		return false
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
	)

	colly.NewCollector()

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

// It scrapes latitude and longitude from parameter 'url' and
// returns them as float64 values
func getLatLongFromURL(url string) (lat, long float64) {
	parts := strings.Split(url, "&")
	latLong := strings.Split(parts[0], "=")[1]
	latStr, longStr, _ := strings.Cut(latLong, ",")

	lat, _ = strconv.ParseFloat(latStr, 64)
	long, _ = strconv.ParseFloat(longStr, 64)
	return
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
	RoomID   string // ID in the database table
	RoomFinderID string
	RoomLoad int
}

// This function prepares the roomIDs that will later be concatenated to the RoomFinder's url/path
// to get the room map and then scrape its coordinate. It returns a slice of all roomIDs
func PrepareDataToScrape() ([]RoomInfo, int) {
	res := DBService.FetchApstatData("apstat")
	fmt.Printf("A number of %d datapoints was fetched.", len(res))
	var data []RoomInfo
	var total int
	for _, val := range res {
		if val.RF_ID != "unknown" {
			continue
		}
		roomNr := scrapeRoomNrFromRoomName(val.Room)
		buildingNr := scrapeBuildingNrFromAddress(val.Address)
		currTotalLoad := getCurrentTotalLoad(val.Load)
		total += currTotalLoad
		roomFinderID := fmt.Sprintf("%s@%s", roomNr, buildingNr)
		data = append(data, RoomInfo{RoomID: val.ID, RoomFinderID: roomFinderID, RoomLoad: currTotalLoad})
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

func UpdateRoomHasLocation(dbName string) {
	db, sqlError := sql.Open("sqlite3", "./data/sqlite/apstat.db")
	if sqlError != nil {
		log.Panicf("Error: could not open sqlite instance! %s", sqlError)
	}
	if db == nil {
		panic("db is nil")
	}
	query := fmt.Sprintf(`UPDATE %s SET RF_ID = ? WHERE RF_ID='unknown'`, dbName)
	stmt, dbError := db.Prepare(query)

	if dbError != nil {
		panic(dbError)
	}

	_, err := stmt.Exec("unknown")
	if err != nil {
		panic(err)
	}
}
