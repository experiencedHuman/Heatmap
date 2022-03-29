package RoomFinder

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

type MockData struct {
	RoomNumber string
	BuildingNr string
}

func Scrape() {
	// url := "http://portal.mytum.de/displayRoomMap?"

	data := []MockData{
		{RoomNumber: "5501", BuildingNr: "5509"},
		{RoomNumber: "5502", BuildingNr: "5508"},
	}

	// Instantiate default collector
	c := colly.NewCollector()

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	for _, val := range data {
		// Add URLs to the queue
		q.AddURL(fmt.Sprintf("http://portal.mytum.de/displayRoomMap?%s@%s", val.RoomNumber, val.BuildingNr))
	}
	// Consume URLs
	q.Run(c)
}
