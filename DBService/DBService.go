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
	forecastTable = "forecast"
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

func RetrieveAccessPointByName(db *sql.DB, name string) *AccessPoint {
	stmt := fmt.Sprintf(`
		SELECT Name, Lat, Long, Load
		FROM apstat
		WHERE Name='%s'
	`, name)

	row := db.QueryRow(stmt)
	result := AccessPoint{}
	switch err := row.Scan(&result.Name, &result.Lat, &result.Long, &result.Load); err {
	case nil:
		log.Printf("Returning result for access point with name %s!", name)
	default:
		log.Printf("No access point found with name %s in the DB!", name)
		log.Println("Returning empty result.")
	}
	return &result
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

// Stores csv data in SQLite table
// under path parameter 'dbPath'.
func StoreApstatInDB(csvPath, dbPath, tableName string) {
	csvData := readFromCSV(csvPath)
	db := InitDB(dbPath)

	createTable := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS "%s" (
				"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"Address"	TEXT,
				"Room"		TEXT,
				"Name"		TEXT,
				"Floor"		TEXT,
				"Status"	TEXT,
				"Type"		TEXT,
				"Load" 		TEXT,
				"Lat"		TEXT,
				"Long"		TEXT
			);
		`, tableName)

	runQuery(createTable)

	query := fmt.Sprintf(`
		INSERT INTO %s (Address, Room, Name, Floor, Status, Type, Load) 
		values (?,?,?,?,?,?,?)
	`, tableName)
	stmt, dbError := db.Prepare(query)

	if dbError != nil {
		panic(dbError)
	}

	for r := range csvData {
		address := csvData[r][0]
		room := csvData[r][1]
		apName := csvData[r][2]
		floor := string(apName[6])
		status := csvData[r][3]
		apType := csvData[r][4]
		apLoad := csvData[r][5]

		_, err := stmt.Exec(address, room, apName, floor, status, apType, apLoad)
		if err != nil {
			panic(err)
		}
	}
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
