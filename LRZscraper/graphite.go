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
	Target     string      `json:"target"`
	Datapoints []DataPoint `json:"datapoints"`
}

type DataPoint struct {
	devices   string //nr of connected devices
	timestamp int
}

func GetGraphiteDataAsJSON(apName string, time string) []NetworkLoad {
	if !strings.HasPrefix(apName, "apa") {
		log.Fatalf("Name of the Access Point must start with \"apa\"!")
	}

	url := buildURL(apName, time, "json")
	log.Printf("url: %s", url)
	resp, httpError := http.Get(url)
	if httpError != nil {
		log.Fatalf("Could not retrieve json data from URL! %q", httpError)
	}
	defer resp.Body.Close()

	var networks []NetworkLoad
	err := json.NewDecoder(resp.Body).Decode(&networks)
	if err != nil {
		log.Fatalf("Could not decode JSON response: %v", err)
	}

	for _, network := range networks {
		fmt.Println(network.Target)
	}

	return networks
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
	fmt.Printf("Nr or retrieved APs: %d", len(APs))

	var gesamt []DataPoint
	var totDevices = 0
	for _, ap := range APs {
		networks := GetGraphiteDataAsJSON(ap.Name, "")
		for i := range networks[0].Datapoints {
			n0, err := strconv.Atoi(networks[0].Datapoints[i].devices)
			if err == nil {
				totDevices += n0
			}
			n1, err := strconv.Atoi(networks[1].Datapoints[i].devices)
			if err == nil {
				totDevices += n1
			}
			n2, err := strconv.Atoi(networks[2].Datapoints[i].devices)
			if err == nil {
				totDevices += n2
			}
			n3, err := strconv.Atoi(networks[3].Datapoints[i].devices)
			if err == nil {
				totDevices += n3
			}
			n4, err := strconv.Atoi(networks[4].Datapoints[i].devices)
			if err == nil {
				totDevices += n4
			}

			ts := networks[0].Datapoints[i].timestamp
			newDataPoint := DataPoint{"10", ts}
			gesamt = append(gesamt, newDataPoint)
		}
	}

	storeHourlyAvgForLast30Days(gesamt)
}

func storeHourlyAvgForLast30Days(datapoints []DataPoint) {
	var lastHour *int
	var n, cnt = 0, 0

	for _, datapoint := range datapoints {
		devices, err := strconv.Atoi(datapoint.devices)
		if err == nil {
			n += devices
		}
		t := getTimeFromTimestamp(datapoint.timestamp)
		hr := t.Hour()
		if lastHour == nil {
			*lastHour = hr
		}

		if hr != *lastHour {
			*lastHour = hr
			avg := n / cnt
			fmt.Println(avg) // save avg value for this hr
			// day := t.Day() save day as well in the entry ?
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
