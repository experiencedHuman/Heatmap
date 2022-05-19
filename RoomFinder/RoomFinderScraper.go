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

var (
	validReq       = 0
	failedRequests = 0
)

type Result struct {
	Failures  []Failure
	Successes []Success
}

type Failure struct {
	ID             string // primary key of the AP
	RoomNr         string
	BuildingNr     string
}

type Success struct {
	ID   string // primary key of the AP
	Lat  string
	Long string
}

type RF_Info struct {
	ID           string // primary key in 'apstat' table
	roomNr       string // ~ architectNr
	buildingNr   string
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

		data = append(data,
			RF_Info{ap.ID,
				architectNr,
				buildingNr,
				roomFinderID,
				currTotalLoad,
				url})
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
	} else if buildingNr == "" {
		if address == "TUM, Tentomax Weihenstephan (Prüfungszelt)Gregor-Mendel-Str.Freising" {
			buildingNr = "4298"
		} else {
			fmt.Println("BuildingNr is empty:", address)
		}
	} else if buildingNr =="8102" {
		fmt.Println("Address of 8102:", address)
	}
	return buildingNr
}

// This function scrapes the roomNr from a longer room description and returns it
func scrapeRoomNrFromRoomName(roomName string) string {
	re := regexp.MustCompile("[0-9]+.[0-9]+(.[0-9])?")
	roomNr := re.FindString(roomName)

	if strings.Contains(roomNr, " ") {
		nums := strings.Split(roomNr, " ")
		
		if len(nums[0]) > len(nums[1]) {
			if nums[1] == "0" {
				roomNr = fmt.Sprintf("EG.%s", nums[0])
			} else {
				roomNr = fmt.Sprintf("0%s.%s", nums[1], nums[0])
			}

		} else if len(nums[0]) == 1 && len(nums[1]) == 2 {
			if nums[0] == "0" {
				roomNr = fmt.Sprintf("EG.0%s", nums[1])
			} else {
				roomNr = fmt.Sprintf("0%s.0%s", nums[0], nums[1])
			}

		} else {
			roomNr = fmt.Sprintf("0%s.%s", nums[0], nums[1])
		}

	} else if roomNr == "" {
		if (roomName == "1.OG" || roomName == "2.OG") {
			roomNr = roomName
		} else {
			fmt.Println("RoomNr is empty:", roomName)
		}

	} else if len(roomNr) == 4 {
		roomNr = fmt.Sprintf("0%c.%v", roomNr[0], roomNr[1:4])
	}
	return roomNr
}

func GetCurrentTotalLoad(load string) int {
	// this regex must match a substring beginning with '(', ignores first number and '-', and then gets the second number
	re := regexp.MustCompile(`\(\s*\d+\s*-\s*(\d+)`)
	match := re.FindStringSubmatch(load)
	if len(match) <= 1 {
		log.Println("FIXME")
		return 0
	}
	currentLoad, err := strconv.Atoi(match[1]) // TODO error handling
	if err != nil {
		log.Println(err)
		return 0 // TODO handle edge cases
	} else {
		return currentLoad
	}
}

func ScrapeURLs(rfInfos []RF_Info) Result {
	var result Result
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
		go scrapeURL(rfInfo, &wg, &t, &result)
	}

	wg.Wait()
	elapsed := time.Since(start)

	totalRes := len(rfInfos)
	log.Printf("Failed requests: %d out of %d\n", failedRequests, totalRes)
	log.Printf("Valid requests: %d out of %d", validReq, totalRes)

	log.Println("time elapsed:", elapsed)

	return result
}

func scrapeURL(rfInfo RF_Info, wg *sync.WaitGroup, t *http.Transport, result *Result) {
	defer wg.Done()
	c := &http.Client{
		Transport: t,
	}
	resp, err := c.Get(rfInfo.url)

	if err != nil {
		log.Println("Failed request url:", rfInfo.url) //TODO uncomment
		failedRequests++
		// TODO collect failed URLs and try NavigaTUM
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Failed request url:", rfInfo.url) //TODO uncomment
		failedRequests++
		// TODO collect failed URLs and try NavigaTUM
		return
	}

	// retrieve HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Println("Failed request url:", rfInfo.url)
		log.Fatal(err)
	}

	// check for link with google map's coordinate
	element := doc.Find("a[href^='http://maps.google.com']")
	link, exists := element.Attr("href")

	if exists {
		validReq++
		lat, long := getLatLongFromURL(link)
		success := Success{rfInfo.ID, lat, long}
		result.Successes = append(result.Successes, success)
	} else {
		// log.Println("Link doesnt exist:", rfInfo.url)
		failure := Failure{rfInfo.ID, rfInfo.roomNr, rfInfo.buildingNr}
		result.Failures = append(result.Failures, failure)
	}
}

// It scrapes latitude and longitude from url and
// returns them as strings
func getLatLongFromURL(url string) (lat, long string) {
	parts := strings.Split(url, "&")
	latLong := strings.Split(parts[0], "=")[1]
	lat, long, _ = strings.Cut(latLong, ",")
	return
}
