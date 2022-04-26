package DBService

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type AccessPoint struct {
	ID      string	// primary key
	Address string
	Room    string
	Name    string
	Floor	string
	Status  string
	Type    string
	Load    string
	Lat     string
	Long    string
}

type APLoad struct {
	Name    string	// primary key
	Network string
	Current string
	Max     string
	Min     string
	Avg     string
}

// Initializes a local database instance, located in dbPath
// returns a pointer to the initialized database.
func InitDB(dbPath string) *sql.DB {
	db, sqlError := sql.Open("sqlite3", dbPath)
	if sqlError != nil {
		log.Panicf("Could not open sqlite instance at path %s! %s", dbPath, sqlError)
	}
	if db == nil {
		panic("db is nil")
	}
	return db
}

// Queries 'accesspoints' table and 
// returns rows where primary key Name = 'name'
func RetrieveAPLoads(db *sql.DB, name string) []APLoad {
	query := fmt.Sprintf(`
		SELECT Name, Network, Current, Max, Min, Avg
		FROM accesspoints
		WHERE Name='%s'
	`, name)
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []APLoad
	for rows.Next() {
		item := APLoad{}
		err2 := rows.Scan(&item.Name, 
						  &item.Network, 
						  &item.Current, 
						  &item.Max, 
						  &item.Min, 
						  &item.Avg)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	return result
}

// Queries 'apstat' table and
// returns all rows where 'address' contains "TUM" and 
// Lat, Long are unassigned
func RetrieveAPs(db *sql.DB, withCoordinate bool) []AccessPoint {
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

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []AccessPoint
	for rows.Next() {
		item := AccessPoint{}
		err2 := rows.Scan(&item.ID, 
						  &item.Address, 
						  &item.Room,
						  &item.Name,
						  &item.Floor, 
						  &item.Load,
						  &item.Lat, 
						  &item.Long)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	return result
}

func StoreAPLoadInDB(csvPath, dbPath, tableName string) {
	csvData := readFromCSV(csvPath)
	db := InitDB(dbPath)

	createTable := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS "%s" (
				"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"network"	TEXT,
				"current"	TEXT,
				"max"		TEXT,
				"min"		TEXT,
				"avg" 		TEXT
			);
		`, tableName)

	runQuery(db, createTable)

	stmt, dbError := db.Prepare(`
		INSERT INTO accesspoints (network, current, max, min, other) values (?,?,?,?,?)
	`)

	if dbError != nil {
		panic(dbError)
	}

	for r := range csvData {
		network := csvData[r][0]
		current := csvData[r][1]
		max 	:= csvData[r][2]
		min 	:= csvData[r][3]
		avg 	:= csvData[r][4]

		_, err := stmt.Exec(network, current, max, min, avg)
		if err != nil {
			panic(err)
		}
	}
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

	runQuery(db, createTable)

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
		room 	:= csvData[r][1]
		apName  := csvData[r][2]
		floor 	:= string(apName[6])
		status  := csvData[r][3]
		apType  := csvData[r][4]
		apLoad  := csvData[r][5]

		_, err := stmt.Exec(address, room, apName, floor, status, apType, apLoad)
		if err != nil {
			panic(err)
		}
	}
}

// Adds a new column of type 'TEXT' to table 'tableName'.
func AddColumn(db *sql.DB, tableName, newCol string) {
	query := fmt.Sprintf(`ALTER TABLE %s ADD COLUMN %s TEXT`, tableName, newCol)
	var args []interface{}
	runQuery(db, query, args...)
}

// Updates 'column' at rows satisfying the 'where' condition with 'newValue'.
func UpdateColumn(db *sql.DB, tableName, column, newValue, where string) {
	query := fmt.Sprintf(`UPDATE %s SET %s = ? WHERE %s`, tableName, column, where)
	runQuery(db, query, newValue)
}

// Updates name of 'currName' column to 'newName'.
func UpdateColumnName(db *sql.DB, tableName, currName, newName string) {
	query := fmt.Sprintf(`ALTER TABLE %s RENAME COLUMN %s TO %s;`, tableName, currName, newName)
	runQuery(db, query)
}

// Prepares sqlite 'query' and executes it with optional 'params'.
func runQuery(db *sql.DB, query string, params ...interface{}) {
	stmt, err := db.Prepare(query)

	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	
	_, err = stmt.Exec(params...)
	if err != nil {
		panic(err)
	}
}
