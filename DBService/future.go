package DBService

import (
	"fmt"
)

func CreateFutureTable(tableName string) {
	hourColumns := createColumnsFor24Hrs()
	createQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS "%s" (
		ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		AP_Name TEXT NOT NULL,
		Day INTEGER,
		%s
	);
	`, tableName, hourColumns)
	db := InitDB(heatmapDB)
	runQuery(createQuery)
	db.Close()
}

// Populates the future table with 15 name entries
// for each access point (one entry per day).
func PopulateFutureTable(accessPoints []AccessPoint) {
	db := InitDB(heatmapDB)
	for _, ap := range accessPoints {
		for j := 0; j < 15; j++ {
			InsertFuture(ap.Name, j)
		}
	}
	db.Close()
}

func InsertFuture(apName string, day int) {
	insertQuery := "INSERT INTO future(AP_Name, Day) VALUES (?, ?)"
	runQuery(insertQuery, apName, day)
}