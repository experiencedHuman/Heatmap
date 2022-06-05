package LRZscraper

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kvogli/Heatmap/DBService"
)

const (
	ApstatTable = "./data/sqlite/apstat.db"
	ApstatCSV   = "data/csv/apstat.csv"
)

var apstatDB = DBService.InitDB(ApstatTable)

const (
	dstPath = "data/csv/"

	rendererEndpoint = "http://graphite-kom.srv.lrz.de/render/"
	// title            = "SSIDs%20(weekly)" //irrelevant
	xFormat  = "%25d.%25m."
	timezone = "CET"
	fontName = "&fontName=Courier"
	title    = "&title=SSIDs%20(weekly)"
)

func getNetworkAlias(apName, network string) string {
	return fmt.Sprintf("alias(ap.%s.ssid.%s,%%22%s%%22)", apName, network, network)
}

func buildURL(apName string, time string, format string) (url string) {
	width := "640"
	height := "240"
	areaMode := "stacked"
	from := "-30days" // todo implement it to use time
	eduroam := getNetworkAlias(apName, "eduroam")
	lrz := getNetworkAlias(apName, "lrz")
	mwn_evt := getNetworkAlias(apName, "mwn-events")
	bayWLAN := getNetworkAlias(apName, "@BayernWLAN")
	other := getNetworkAlias(apName, "other")
	target := fmt.Sprintf("cactiStyle(group(%s,%s,%s,%s,%s))", eduroam, lrz, mwn_evt, bayWLAN, other)
	url = fmt.Sprintf(`%s?width=%s&height=%s&title=%s&areaMode=%s&xFormat=%s&tz=%s&from=%s&target=%s&fontName=Courier&format=%s`, rendererEndpoint, width, height, title, areaMode, xFormat, timezone, from, target, format)
	return
}

// makes a GET request to LRZ's Graphite "/renderer" endpoint
// and stores the retrieved data in data/csv/{apName}.csv
func GetGraphiteDataAsCSV(apName string, time string) {
	if !strings.HasPrefix(apName, "apa") {
		log.Fatalf("Name of the Access Point must start with \"apa\"!")
	}

	url := buildURL(apName, time, "csv")
	resp, httpError := http.Get(url)
	if httpError != nil {
		log.Fatalf("Could not retrieve csv data from URL! %q", httpError)
	}

	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)
	csvReader.Comma = ','

	dst := fmt.Sprintf("%s%s.csv", dstPath, apName)
	file, osError := os.Create(dst)
	if osError != nil {
		log.Fatalf("Could not create file, err: %q", osError)
	}
	defer file.Close()

	// TODO clean duplicate code
	dst = fmt.Sprintf("%sprophet-%s.csv", dstPath, apName)
	prophetFile, osError := os.Create(dst)
	if osError != nil {
		log.Fatalf("Could not create file, err: %q", osError)
	}
	defer prophetFile.Close()

	responseCSV := csv.NewWriter(file)
	responseCSV.Comma = ';'
	defer responseCSV.Flush()

	// TODO clean duplicate code
	prophetCSV := csv.NewWriter(prophetFile)
	prophetCSV.Comma = ','
	defer prophetCSV.Flush()
	prophetCSV.Write([]string{
		"ds", "y",
	})

	var csvRecord []string
	for i := 0; ; i++ {
		csvRecord, httpError = csvReader.Read()
		if httpError == io.EOF {
			break
		} else if httpError != nil {
			panic(httpError)
		} else {
			fmt.Println(csvRecord)
			fields := strings.Fields(csvRecord[0]) // get substrings separated by whitespaces
			network := fields[0]
			current := strings.Split(fields[1], ":")[1]
			max := strings.Split(fields[2], ":")[1]
			min := strings.Split(fields[3], ":")[1]

			dateAndTime := csvRecord[1]
			avg := csvRecord[2]

			responseCSV.Write([]string{
				network, current, max, min, dateAndTime, avg,
			})

			// TODO clean duplicate code
			ds := fmt.Sprintf("%s", dateAndTime)
			if current == "nan" {
				current = "0.0"
			}
			prophetCSV.Write([]string{
				ds, current,
			})
		}
	}

	log.Printf("Data was stored in %s", dst)
}

type NetworkLoad struct {
	Datapoints [][]float64 `json:"datapoints"`
	Target     string      `json:"target"`
}

func GetGraphiteDataAsJSON(apName string, time string) ([]NetworkLoad, error) {
	if !strings.HasPrefix(apName, "apa") {
		log.Fatalf("Name of the Access Point must start with \"apa\"!")
	}

	url := buildURL(apName, time, "json")
	log.Printf("url: %s", url)
	resp, httpError := http.Get(url)
	if httpError != nil {
		log.Printf("Could not retrieve json data from URL! %q", httpError)
		return nil, httpError
	}
	defer resp.Body.Close()

	var networks []NetworkLoad
	err := json.NewDecoder(resp.Body).Decode(&networks)
	if err != nil {
		// sometimes server sends bad JSON resp
		log.Printf("Could not decode JSON response: %v", err)
		return nil, err
	}

	// for _, network := range networks {
	// 	fmt.Println(network.Target)
	// }

	return networks, nil
}

func GetMessage(networks []NetworkLoad) string {
	totCurr, totMax, totMin := 0.0, 0.0, 0.0
	for _, network := range networks {
		curr, max, min := getCurrMaxMin(network.Target)
		totCurr += curr
		totMax += max
		totMin += min
	}

	diff := totMax - totMin
	if diff > 0.0 {
		if totCurr > 0.85*totMin {
			return "Queue of death"
		} else if totCurr > 0.6*totMin {
			return "Batch processing"
		} else if totCurr > 0.3*totMin {
			return "Here we go!"
		} else {
			return "All good"
		}
	}
	return "Keine Information!"
}

func getCurrMaxMin(networkLoad string) (curr, max, min float64) {
	fields := strings.Fields(networkLoad) //remove whitespaces

	currStr := strings.Split(fields[0], ":")[1]
	curr, err := strconv.ParseFloat(currStr, 32)
	if err != nil {
		curr = 0.0
	}

	maxStr := strings.Split(fields[1], ":")[1]
	max, err = strconv.ParseFloat(maxStr, 32)
	if err != nil {
		max = 0.0
	}

	minStr := strings.Split(fields[2], ":")[1]
	min, err = strconv.ParseFloat(minStr, 32)
	if err != nil {
		min = 0.0
	}

	return
}

func Last30Days() {
	// get list of 2400 APs
	// get json data for each one for the last 30 days
	// group data by hour of each day
	// save it in last30days.csv (or sqlite table)
	// 2400 aps * 30 days * 24 hr = 1,728,000 entries
	APs := DBService.RetrieveAPsOfTUM(apstatDB, true)
	unprocessedAPs := DBService.GetUnprocessedAPs()
	
	// DBService.InsertAPsLast30Days(APs)
	fmt.Printf("Nr or retrieved APs: %d\n", len(APs))

	for i, ap := range APs {
		_, ok := unprocessedAPs[ap.Name]
		if ok {
			SaveLast30DaysForAP(ap)
			log.Printf("%d of %d done.", i, len(APs))
		}
	}
}

func SaveLast30DaysForAP(ap DBService.AccessPoint) {
	networks, err := GetGraphiteDataAsJSON(ap.Name, "")
	if err != nil {
		fmt.Printf("Skipping %s. Bad JSON response from server.\n", ap.Name)
		return
	}
	gesamt := networks[0].Datapoints

	n := len(networks)
	for i := 1; i < n; i++ {
		printed := false
		for j := range gesamt {
			lenGesamt := len(gesamt)
			lenOtherNetwork := len(networks[i].Datapoints)
			if lenGesamt != lenOtherNetwork && !printed {
				log.Printf("Different lengths: %d vs %d", lenGesamt, lenOtherNetwork)
				printed = true
			}

			if j < lenOtherNetwork {
				gesamt[j][0] += networks[i].Datapoints[j][0]
			}
		}
	}
	storeHourlyAvgForLast30Days(gesamt, ap.Name)
}

// for every hour there are 4 datapoints being collected, one each 15 min
// calculate hourly avg and store it in the database
func storeHourlyAvgForLast30Days(datapoints [][]float64, apName string) {
	var lastHour *int
	var n, cnt = 0, 0

	for _, datapoint := range datapoints {
		devices := datapoint[0]
		n += int(devices)

		ts := int(datapoint[1])
		t := getTimeFromTimestamp(ts)
		hr := t.Hour()
		if lastHour == nil {
			lastHour = &hr
			// fmt.Println(*lastHour)
		}

		if hr != *lastHour {
			*lastHour = hr
			avg := n / cnt
			// save avg value for this hr
			// day := t.Day() save day as well in the entry ?
			// fmt.Printf("hr: %d, avg: %d, day: %d\n", t.Hour(), avg, t.Day())
			DBService.UpdateLast30Days(t.Day(), t.Hour(), avg, apName)
			cnt = 0
			n = 0
		} else {
			cnt += 1
		}
	}
}

// parses a Unix timestamp (i.e. milliseconds from EPOCH) and
// returns it as time.Time
func getTimeFromTimestamp(timestamp int) time.Time {
	ts := fmt.Sprintf("%d", timestamp)
	t, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(t, 0)
	return tm
}
