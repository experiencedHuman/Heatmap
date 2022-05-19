package NavigaTUM

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

// makes an HTTP request to nav.tum.sexy/api/get/<roomID>
// e.g. roomID := "5602.EG.001"
func GetRoomCoordinates(roomID string) (lat, long string, found bool) {
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
