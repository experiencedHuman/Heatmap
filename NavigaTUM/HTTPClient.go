package NavigaTUM

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kvogli/Heatmap/RoomFinder"
)

type CoordInfo struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"lon"`
	Src  string  `json:"source"`
}

type NavigatumResponse struct {
	Coords CoordInfo `json:"coords"`
}

func GetRoomCoordinates(roomID string) RoomFinder.Coordinate {
	// e.g. roomID := "5602.EG.001"
	url := fmt.Sprintf("https://roomapi.tum.sexy/api/get/%s", roomID)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error: GET request to %s failed!", url)
	}
	var naviResp NavigatumResponse
	err = json.NewDecoder(resp.Body).Decode(&naviResp)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("Error: JSON decoding failed!")
	}

	return RoomFinder.Coordinate{Lat: naviResp.Coords.Lat, Long: naviResp.Coords.Long}
}
