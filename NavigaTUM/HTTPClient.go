package NavigaTUM

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"net/http"
	"github.com/kvogli/Heatmap/RoomFinder"
	"github.com/kvogli/Heatmap/DBService"
)

const (
	ApstatTable = "./data/sqlite/apstat.db"
)

type Coords struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"lon"`
	Src  string  `json:"source"`
}

// Navigatum Response
type NaviRes struct {
	Type  string `json:"type"`
	Coord Coords `json:"coords"`
}

func ScrapeNavigaTUM(res RoomFinder.Result) (count int) {
	count = 0 // number of found coordinates
	db := DBService.InitDB(ApstatTable)

	for _, res := range res.Failures {
		var roomID string
		if strings.Contains(res.RoomNr, "OG") || res.RoomNr == "" || strings.Contains(res.RoomNr, "..") {
			roomID = res.BuildingNr
		} else {
			roomID = fmt.Sprintf("%s.%s", res.BuildingNr, res.RoomNr)
		}

		lat, long, found := getRoomCoordinates(roomID)

		if found {
			where := fmt.Sprintf("ID='%s'", res.ID)
			DBService.UpdateColumn(db, "apstat", "Lat", lat, where)
			DBService.UpdateColumn(db, "apstat", "Long", long, where)
			count++
		} else {
			lat, long, found = getRoomCoordinates(res.BuildingNr)
			if found {
				where := fmt.Sprintf("ID='%s'", res.ID)
				DBService.UpdateColumn(db, "apstat", "Lat", lat, where)
				DBService.UpdateColumn(db, "apstat", "Long", long, where)
				count++
			}
		}
	}

	return
}

// makes an HTTP GET request to nav.tum.sexy/api/get/{roomID}
// e.g. roomID := 5602.EG.001
func getRoomCoordinates(roomID string) (lat, long string, found bool) {
	lat, long, found = "", "", false

	url := fmt.Sprintf("https://nav.tum.sexy/api/get/%s", roomID)
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("GET request to %s failed!", url)
		return
	}

	if resp.StatusCode != 200 {
		log.Printf("%v for url: %s", resp.Status, url)
		return
	}

	var nResp NaviRes
	err = json.NewDecoder(resp.Body).Decode(&nResp)
	defer resp.Body.Close()

	if err != nil {
		log.Printf("JSON decoding failed! %q", err)
		log.Printf("%v", resp.Body)
		return
	}

	lat = fmt.Sprintf("%f", nResp.Coord.Lat)
	long = fmt.Sprintf("%f", nResp.Coord.Long)
	found = true

	typ := nResp.Type
	fmt.Println(typ)

	return
}
