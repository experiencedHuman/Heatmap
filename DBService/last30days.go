package DBService

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateHistoryTable(tableName string, dbName string) {
	hourColumns := createColumnsFor24Hrs()
	createQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS "%s" (
		ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		AP_Name TEXT NOT NULL,
		Day INTEGER,
		%s
	);
	`, tableName, hourColumns)
	db := InitDB(dbName)
	runQuery(db, createQuery)
	db.Close()
}


func createColumnsFor24Hrs() string {
	partOfQuery := "T0 INTEGER"
	for i := 1; i < 24; i++ {
		partOfQuery = fmt.Sprintf("%s,\nT%d INTEGER", partOfQuery, i)
	}
	return partOfQuery
}

func UpdateLast30Days(day int, hour int, avg int, apName string) {
	updateQuery := fmt.Sprintf(`
		UPDATE last30days SET Day = ?, T%d = ? WHERE AP_NAME = '%s'
	`, hour, apName)
	db := InitDB(last30daysTable)
	runQuery(db, updateQuery, day, avg)
	db.Close()
}

func UpdateLast31Days(day int, hour int, avg int, apName string) {
	updateQuery := fmt.Sprintf(`
		UPDATE last31days SET Day = ?, T%d = ? WHERE AP_NAME = '%s'
	`, hour, apName)
	db := InitDB(last30daysTable)
	runQuery(db, updateQuery, day, avg)
	db.Close()
}

// Populates the last31days table with 31 name entries
// for each access point (one entry per day).
func PopulateLast30DaysTable(accessPoints []AccessPoint) {
	db := InitDB(historyDB)
	for _, ap := range accessPoints {
		for j := 0; j < 31; j++ {
			InsertLast30Days(ap.Name, db)
			fmt.Println("Inserted")
		}
	}
	db.Close()
}

func InsertLast30Days(apName string, db *sql.DB) {
	insertQuery := "INSERT INTO last31days(AP_Name) VALUES (?)"
	runQuery(db, insertQuery, apName)
}

// returns a map of AP names, whose last 30 day data
// hasn't been stored in the database yet
func GetUnprocessedAPs() map[string]bool {
	db := InitDB(last30daysTable)
	query := `
		SELECT DISTINCT AP_Name
		FROM last30days
		WHERE Day IS NULL
	`
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
	db := InitDB(last30daysTable)
	query := fmt.Sprintf(`
		SELECT T%d
		FROM last30days
		WHERE AP_Name = '%s'
		AND Day = %d
	`, hr, name, day)

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
	db := InitDB(last30daysTable)
	query := fmt.Sprintf(`
		SELECT DISTINCT AP_Name, T%d
		FROM last30days
		WHERE Day = %d
	`, hr, day)
	
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []AccessPoint
	result = make([]AccessPoint, 0)
	
	for rows.Next() {
		var ap AccessPoint
		err = rows.Scan(&ap.Name, &ap.Load);
		if err == nil {
			result = append(result, ap)
		} else {
			panic(err)
		}
	}
	
	db.Close()
	return result
}