package DBService

import (
	"fmt"
	"log"
)

func CreateFutureTable() {
	hourColumns := createColumnsFor24Hrs()
	createQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS "future" (
		ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		AP_Name TEXT NOT NULL,
		Day INTEGER,
		%s
	);
	`, hourColumns)
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

func UpdateTomorrow() {
	deleteDayEq30 := "DELETE FROM future WHERE Day = 0"
	runQuery(deleteDayEq30)
	
	decreaseDays := "UPDATE future SET Day = Day - 1 "
	runQuery(decreaseDays)
}

func GetPredictionForSingleAP(day int, hour int, name string) AccessPoint {
	db := InitDB(heatmapDB)
	query := fmt.Sprintf(`
		SELECT T%d
		FROM future
		WHERE AP_Name = '%s'
		AND Day = %d
	`, hour, name, day)

	row := db.QueryRow(query)
	result := AccessPoint{}
	switch err := row.Scan(&result.Load); err {
	case nil:
		log.Println("Returning result from history!")
	default:
		log.Println("No data found in history! Returning empty result.")
	}
	
	db.Close()
	return result
}

func GetFutureForAllAPs(day int, hr int) []AccessPoint {
	db := InitDB(heatmapDB)
	query := fmt.Sprintf(`
		SELECT DISTINCT AP_Name, T%d, Max, Min
		FROM future
		WHERE Day = %d
	`, hr, day)
	
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var apList []AccessPoint
	apList = make([]AccessPoint, 0)
	
	for rows.Next() {
		var ap AccessPoint
		err = rows.Scan(&ap.Name, &ap.Load, &ap.Max, &ap.Min);
		if err == nil {
			apList = append(apList, ap)
		} else {
			// log.Println(err) // TODO
		}
	}
	
	db.Close()
	return apList
}