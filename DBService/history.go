package DBService

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateHistoryTable(tableName string) {
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


func createColumnsFor24Hrs() string {
	partOfQuery := "T0 INTEGER"
	for i := 1; i < 24; i++ {
		partOfQuery = fmt.Sprintf("%s,\nT%d INTEGER", partOfQuery, i)
	}
	return partOfQuery
}

func UpdateHistory(day int, hour int, avg int, apName string) {
	updateQuery := fmt.Sprintf(`
		UPDATE history SET T%d = ? WHERE AP_NAME = '%s' AND Day = '%d'
	`, hour, apName, day)
	runQuery(updateQuery, avg)
}

// Populates the last31days table with 31 name entries
// for each access point (one entry per day).
func PopulateHistoryTable(accessPoints []AccessPoint) {
	db := InitDB(heatmapDB)
	for _, ap := range accessPoints {
		for j := 0; j < 31; j++ {
			InsertHistory(ap.Name, j, db)
		}
	}
	db.Close()
}

func InsertHistory(apName string, day int, db *sql.DB) {
	insertQuery := "INSERT INTO history(AP_Name, Day) VALUES (?, ?)"
	runQuery(insertQuery, apName, day)
}

// returns a map of AP names, whose last 30 day data
// hasn't been stored in the database yet
func GetUnprocessedAPs() map[string]bool {
	db := InitDB(heatmapDB)
	query := fmt.Sprintf(`
		SELECT DISTINCT AP_Name
		FROM %s
		WHERE T0 IS NULL
	`, historyTable)

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return make(map[string]bool)
	}
	defer rows.Close()

	var names = make(map[string]bool)
	for rows.Next() {
		var ap string
		rows.Scan(&ap)
		names[ap] = true
	}

	db.Close()
	return names
}

// queries the database and returns the intensity of
// the access point, based on the selected day and hour
func GetHistoryOfSingleAP(name string, day int, hr int) AccessPoint {
	db := InitDB(heatmapDB)
	query := fmt.Sprintf(`
		SELECT T%d
		FROM %s
		WHERE AP_Name = '%s'
		AND Day = %d
	`, hr, historyTable, name, day)

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

// queries the database and returns a list of the names and intensities (network load)
// of all access points, based on the selected day and hour.
func GetHistoryOfAllAccessPoints(day int, hr int) []AccessPoint {
	db := InitDB(heatmapDB)
	query := fmt.Sprintf(`
		SELECT DISTINCT AP_Name, T%d, Max, Min
		FROM %s
		WHERE Day = %d
	`, hr, historyTable, day)
	
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