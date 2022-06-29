package DBService

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// database path
	heatmapDB = "./data/sqlite/heatmap.db"

	// tables inside heatmapDB
	apstatTable   = "apstat"
	historyTable  = "history"
	forecastTable = "future"
)

var DB = InitDB(heatmapDB)

type AccessPoint struct {
	ID      string // primary key
	Address string
	Room    string
	Name    string
	Floor   string
	Status  string
	Type    string
	Load    string
	Lat     string
	Long    string
	Max			int
	Min			int
}

// Opens a database at dbPath and
// returns a pointer to it.
func InitDB(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Panicf("Could not open SQLite DB! %v", err)
	}

	if db == nil {
		panic("DB is nil")
	}
	return db
}

func RetrieveAPsOfTUM(withCoordinate bool) []AccessPoint {
	var query string

	if withCoordinate {
		query = `
			SELECT ID, Address, Room, Name, Floor, Load, Lat, Long
			FROM apstat
			WHERE Address LIKE '%TUM%'
			AND Lat!='lat'
			AND Long!='long'
		`
	} else {
		query = `
			SELECT ID, Address, Room, Name, Floor, Load, Lat, Long
			FROM apstat
			WHERE Address LIKE '%TUM%'
			AND Lat='lat'
			AND Long='long'
		`
	}

	return RetrieveAPsFromTUM(query)
}

// Queries 'apstat' table and
// returns all rows where 'address' contains "TUM" and
// Lat, Long are unassigned
func RetrieveAPsFromTUM(query string) []AccessPoint {
	db := InitDB(heatmapDB)

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []AccessPoint
	for rows.Next() {
		item := AccessPoint{}
		err := rows.Scan(&item.ID, &item.Address, &item.Room, &item.Name, &item.Floor, &item.Load, &item.Lat, &item.Long)

		if err != nil {
			panic(err)
		}
		result = append(result, item)
	}

	db.Close()
	return result
}

// Reads data from a csv file and
// returns it as a string matrix.
func readFromCSV(csvPath string) [][]string {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'
	csvData, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return csvData
}

// Adds a new column of type 'TEXT' to table 'tableName'.
func AddColumn(tableName, newCol, colType string) {
	query := fmt.Sprintf(`ALTER TABLE %s ADD COLUMN %s %s`, tableName, newCol, colType)
	var args []interface{}
	runQuery(query, args...)
}

// Updates 'column' at rows satisfying the 'where' condition with 'newValue'.
func UpdateColumn(tableName, column, newValue, where string) {
	query := fmt.Sprintf(`UPDATE %s SET %s = ? WHERE %s`, tableName, column, where)
	runQuery(query, newValue)
}

func UpdateColumnInt(tableName, column string, newValue int, where string) {
	query := fmt.Sprintf(`UPDATE %s SET %s = ? WHERE %s`, tableName, column, where)
	runQuery(query, newValue)
}

// Updates name of 'currName' column to 'newName'.
func UpdateColumnName(tableName, currName, newName string) {
	query := fmt.Sprintf(`ALTER TABLE %s RENAME COLUMN %s TO %s;`, tableName, currName, newName)
	runQuery(query)
}

// Prepares sqlite 'query' and executes it with optional 'params'.
func runQuery(query string, params ...interface{}) {
	stmt, err := DB.Prepare(query)

	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		panic(err)
	}
}
