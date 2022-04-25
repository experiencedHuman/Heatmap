package RoomFinder

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kvogli/Heatmap/DBService"

	_ "github.com/mattn/go-sqlite3"

	"net"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type muxCache struct {
	mux sync.Mutex
	APs map[string]Coordinate
}

var (
	cache          = muxCache{APs: make(map[string]Coordinate)}
	failedRequests = 0
	validReq       = 0
)

type Coordinate struct {
	id   string // primary key of the AP
	Lat  string
	Long string
}

type RF_Info struct {
	id           string // primary key in 'apstat' table
	RoomFinderID string // = <architectNr>@<buildingNr>
	ApLoad       int    // current total load of the AP
	url          string // http://portal.mytum.de/displayRoomMap?<roomFinderID>
}

// It receives as input an array of APs and generates a roomFinderID & URL for each element.
// Returns a slice of RF_Infos, containing RoomFinder URLs.
func PrepareDataToScrape(APs []DBService.AccessPoint) ([]RF_Info, int) {
	var data []RF_Info
	var total int
	for _, ap := range APs {
		architectNr := scrapeRoomNrFromRoomName(ap.Room)
		buildingNr := scrapeBuildingNrFromAddress(ap.Address)
		currTotalLoad := GetCurrentTotalLoad(ap.Load)

		total += currTotalLoad
		roomFinderID := fmt.Sprintf("%s@%s", architectNr, buildingNr)
		url := fmt.Sprintf("http://portal.mytum.de/displayRoomMap?%s", roomFinderID)

		data = append(data, RF_Info{ap.ID, roomFinderID, currTotalLoad, url})
	}
	return data, total
}

// This function scrapes the buildingNr from the address description and returns it
func scrapeBuildingNrFromAddress(address string) string {
	re := regexp.MustCompile("[0-9]+")
	buildingNr := re.FindString(address)
	if buildingNr == "5500" {
		re = regexp.MustCompile("\\d{4}")
		buildingNr = re.FindAllString(address, -1)[1]
	}
	return buildingNr
}

// This function scrapes the roomNr from a longer room description and returns it
func scrapeRoomNrFromRoomName(roomName string) string {
	re := regexp.MustCompile("[0-9]+.[0-9]+(.[0-9])?")
	// TODO measure success rate of this regex
	roomNr := re.FindString(roomName)
	return roomNr
}

func GetCurrentTotalLoad(load string) int {
	// this regex must match a substring beginning with '(', ignores first number and '-', and then gets the second number
	re := regexp.MustCompile(`\(\s*\d+\s*-\s*(\d+)`)
	match := re.FindStringSubmatch(load)
	currentLoad, err := strconv.Atoi(match[1]) // TODO error handling
	if err != nil {
		log.Println(err)
		return 0 // TODO handle edge cases
	} else {
		return currentLoad
	}
}

func ScrapeURLs(rfInfos []RF_Info) map[string]Coordinate {
	var wg sync.WaitGroup

	start := time.Now()
	t := http.Transport{
		Dial: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 60 * time.Second,
		MaxConnsPerHost:     50,
		MaxIdleConns:        50,
	}

	for _, rfInfo := range rfInfos {
		wg.Add(1)
		go scrapeURL(rfInfo, &wg, &t)
	}

	wg.Wait()
	elapsed := time.Since(start)

	totalRes := len(rfInfos)
	log.Printf("Failed requests: %d out of %d\n"+
		"Valid requests: %d out of %d", failedRequests, totalRes, validReq, totalRes)
	
	log.Println("time elapsed:", elapsed)

	return cache.APs
}

func scrapeURL(rfInfo RF_Info, wg *sync.WaitGroup, t *http.Transport) {
	defer wg.Done()
	c := &http.Client{
		Transport: t,
	}
	resp, err := c.Get(rfInfo.url)

	if err != nil {
		failedRequests++
		// TODO collect failed URLs and try NavigaTUM
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		failedRequests++
		// TODO collect failed URLs and try NavigaTUM
		return
	}

	// retrieve HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	validReq++
	element := doc.Find("a[href^='http://maps.google.com']")
	link, exists := element.Attr("href")

	if exists {
		lat, long := getLatLongFromURL(link)
		var coord Coordinate
		// RoomFinder loads a message element in case the location is not that of the requested Room
		// but that of the building (in which that room is located)
		if doc.Find(".message").Length() == 0 {
			coord = Coordinate{rfInfo.id, lat, long}
			cache.storeCoord(rfInfo.id, coord)
		} else {
			coord = Coordinate{rfInfo.id, lat, long}
			cache.storeCoord(rfInfo.id, coord)
		}
	} else {
		log.Println("Link doesnt exist:", rfInfo.url)
	}
}

// stores the coordinate of the Access Point with primary key id in a map
func (cache *muxCache) storeCoord(id string, coord Coordinate) {
	cache.mux.Lock()
	defer cache.mux.Unlock()
	cache.APs[id] = coord
}

// It scrapes latitude and longitude from parameter 'url' and
// returns them as float64 values
func getLatLongFromURL(url string) (lat, long string) {
	parts := strings.Split(url, "&")
	latLong := strings.Split(parts[0], "=")[1]
	lat, long, _ = strings.Cut(latLong, ",")
	return
}
