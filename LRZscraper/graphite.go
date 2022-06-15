package LRZscraper

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kvogli/Heatmap/DBService"
)

const (
	ApstatTable = "./data/sqlite/apstat.db"
	ApstatCSV   = "data/csv/apstat.csv"
	historyDB   = "./data/sqlite/history.db"
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

	responseCSV := csv.NewWriter(file)
	responseCSV.Comma = ';'
	defer responseCSV.Flush()

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
		}
	}

	log.Printf("Data is stored in %s", dst)
}

type NetworkLoad struct {
	Datapoints [][]float64 `json:"datapoints"`
	Target     string      `json:"target"`
}

func GetGraphiteDataAsJSON(apName string, time string, t *http.Transport) ([]NetworkLoad, error) {
	if !strings.HasPrefix(apName, "apa") {
		log.Fatalf("Name of the Access Point must start with \"apa\"!")
	}

	url := buildURL(apName, time, "json")

	c := &http.Client{
		Transport: t,
	}

	resp, httpError := c.Get(url)
	if httpError != nil {
		log.Printf("Could not retrieve json data from URL! %q", httpError)
		return nil, httpError
	}
	defer resp.Body.Close()

	var networks []NetworkLoad
	err := json.NewDecoder(resp.Body).Decode(&networks)
	if err != nil {
		// sometimes server sends bad JSON response
		log.Printf("Could not decode JSON response: %v", err)
		return nil, err
	}

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

type APHistory struct {
	name       string
	historyPtr *history
}

func StoreHistory() {
	APs := DBService.RetrieveAPsOfTUM(true)
	var wg sync.WaitGroup
	channel := make(chan APHistory)

	t := http.Transport{
		Dial: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 60 * time.Second,
		MaxConnsPerHost:     50,
		MaxIdleConns:        50,
	}

	start := time.Now()
	log.Println("Start time:", start)
	wg.Add(len(APs))
	for i, ap := range APs {
		go func(apName string, idx int) {
			SaveLast30DaysForAP(apName, channel, &t, idx)
			wg.Done()
		}(ap.Name, i)
	}

	go func() {
		wg.Wait()
		close(channel)
		elapsed := time.Since(start)
		log.Println("Elapsed time:", elapsed)
	}()

	skippedAPs := make([]string, 0)
	histories := make([]APHistory, 0)
	skipped := 0
	for val := range channel {
		if val.historyPtr == nil {
			skipped += 1
			skippedAPs = append(skippedAPs, val.name)
		} else {
			// storeHistoryOfAP(val.name, history)
			// fmt.Println("Storing", val.name)
			histories = append(histories, val)
			
			n := len(histories)
			if n % 100 == 0 {
				log.Println(n)
			}
		}
	}

	log.Println("Started storing in DB:", time.Now())
	storeHistories(histories)
	fmt.Println("Total nr of skipped APs:", skipped)
	saveSkippedToCSV(skippedAPs)
}

// Saves names of skipped Access Points in a csv file.
func saveSkippedToCSV(skippedAPs []string) {
	file, err := os.Create("data/csv/skipped.csv")
	if err != nil {
		log.Println("Failed to store skipped APs!", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	for _, skippedAP := range skippedAPs {
		if err := w.Write([]string{skippedAP}); err != nil {
			log.Println("Failed to store skipped AP in csv file:", skippedAP)
		}
	}
}

func storeHistories(histories []APHistory) {
	for i, apHistory := range histories {
		apName := apHistory.name
		history := *apHistory.historyPtr
		storeHistoryOfAP(apName, history)
		if i % 100 == 0 {
			log.Println("done storing", i, apName)
		}
	}
}

func storeHistoryOfAP(apName string, history history) {
	days := len(history)
	hours := len(history[0])
	if days != 31 || hours != 24 {
		log.Printf("FIXME: Days: %d, Hours: %d", days, hours)
	}
	for day := 0; day < days; day++ {
		for hr := 0; hr < hours; hr++ {
			avg := history[day][hr]
			// DBService.DB.Close() // just testing // TODO remove line
			DBService.UpdateHistory(day, hr, avg, apName)
		}
	}
}

func Last30Days() {
	// get list of 2400 APs
	// get json data for each one for the last 30 days
	// group data by hour of each day
	// save it in history.csv (or sqlite table)
	// 2400 aps * 30 days * 24 hr = 1,728,000 entries
	APs := DBService.RetrieveAPsOfTUM(true)
	unprocessedAPs := DBService.GetUnprocessedAPs()

	// DBService.PopulateHistory(APs)
	fmt.Printf("Nr or retrieved APs: %d\n", len(APs))

	for i, ap := range APs {
		_, ok := unprocessedAPs[ap.Name]
		if ok {
			// SaveLast30DaysForAP(ap)
			log.Printf("%d of %d done.", i, len(APs))
		}
	}
}

func SaveLast30DaysForAP(apName string, c chan APHistory, t *http.Transport, idx int) {
	networks, err := GetGraphiteDataAsJSON(apName, "", t)
	if err != nil {
		fmt.Printf("Skipping %s. Bad JSON response from server.\n", apName)
		c <- APHistory{apName, nil}
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
	history := storeHourlyAvgForHistory(gesamt)
	c <- APHistory{apName, &history}
	// log.Println("Sending", apName, idx)
}

type history [31][24]int

// It calculates hourly averages of network load for the given AP data.
// Returns a 31x24 matrix, where cell (i,j) holds
// the hourly avg. of day i and hour j.
func storeHourlyAvgForHistory(datapoints [][]float64) history {
	var history history //last 30 days + today
	var prevHour, prevDay *int
	day := 31
	cnt := 0
	avg := 0
	n := 0 // network load (Nr of connected devices)

	for _, datapoint := range datapoints {
		n += int(datapoint[0])
		ts := int(datapoint[1])

		t := getTimeFromTimestamp(ts)
		hr := t.Hour()
		currDay := t.Day()

		if prevHour == nil || prevDay == nil {
			prevHour = &hr
			prevDay = &currDay
		}

		if hr != *prevHour {
			*prevHour = hr
			avg = n / cnt
			history[day-1][hr] = avg
			cnt = 0
			n = 0
		} else {
			cnt += 1
		}

		if currDay != *prevDay {
			*prevDay = currDay
			day -= 1
		}
	}

	history[day-1][*prevHour] = avg
	return history
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
