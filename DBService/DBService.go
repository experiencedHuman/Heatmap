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
	sqlitePath string = "./data/sqlite/"
	csvPath    string = "data/csv/"
)

type AccessPointOverview struct {
	ID      string
	Address string
	Room    string
	Name    string
	Status  string
	Type    string
	Load    string
	RF_ID	string // RoomFinderID ~ architect number e.g. 1302@0103
}

// initializes a local database instance, located in dbPath
// returns a pointer to the initialized database
func initDB(dbPath string) *sql.DB {
	db, sqlError := sql.Open("sqlite3", dbPath)
	if sqlError != nil {
		log.Panicf("Error: could not open sqlite instance at path %s! %s", dbPath, sqlError)
	}
	if db == nil {
		panic("db is nil")
	}
	return db
}

func FetchApstatData(dbName string) []AccessPointOverview {
	dbPath := fmt.Sprintf("%s%s.db", sqlitePath, dbName)
	db := initDB(dbPath)
	return readItem(db)
}

func createTableAccesspoints(tableName string, db *sql.DB) {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS "%s" (
			"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"network"	TEXT,
			"current"	TEXT,
			"max"		TEXT,
			"min"		TEXT,
			"other" 	TEXT
		);
	`, tableName)

	stmt, _ := db.Prepare(query)
	_, err := stmt.Exec()
	if err != nil {
		panic(err)
	}
}

func readItem(db *sql.DB) []AccessPointOverview {
	query := `
		SELECT ID, Address, Room, Load, RF_ID
		FROM apstat
		WHERE Address Like '%TUM%'
	`
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []AccessPointOverview
	for rows.Next() {
		item := AccessPointOverview{}
		err2 := rows.Scan(&item.ID, &item.Address, &item.Room, &item.Load, &item.RF_ID)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	return result
}

// for csv/graphData.csv
func StoreDataInSQLite(dbPath string) {
	csvFile, err := os.Open("csv/graphData.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'
	data, err := csvReader.ReadAll() // TODO use csvReader.Read() for big files
	if err != nil {
		log.Fatal(err)
	}

	db := initDB(dbPath)

	createTableAccesspoints("accesspoints", db)

	// fmt.Println(network, current, max, min, other)
	stmt, dbError := db.Prepare(`
		INSERT INTO accesspoints (network, current, max, min, other) values (?,?,?,?,?)
	`)

	if dbError != nil {
		panic(dbError)
	}

	log.Println("Storing graph data in SQLite ...")
	for r := range data {
		network := data[r][0]
		current := data[r][1]
		max := data[r][2]
		min := data[r][3]
		other := data[r][4]

		_, err := stmt.Exec(network, current, max, min, other)
		if err != nil {
			panic(err)
		}
	}
	log.Println("Finished storing graph data.")
}

// It reads the apstat data from the csv file and
// stores it in a SQLite table under path parameter 'dbPath'
func StoreApstatInSQLite(dbName string) {
	filePath := fmt.Sprintf("%s%s.csv", csvPath, dbName)
	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.Comma = ';'
	csvData, err := csvReader.ReadAll() // TODO use csvReader.Read() for big files
	if err != nil {
		log.Fatal(err)
	}

	dbPath := fmt.Sprintf("%s%s.db", sqlitePath, dbName)
	db := initDB(dbPath)
	createTable("apstat", db)

	query := fmt.Sprintf(`
		INSERT INTO %s (Address, Room, Name, Status, Type, Load) 
		values (?,?,?,?,?,?)
	`, dbName)
	stmt, dbError := db.Prepare(query)

	if dbError != nil {
		panic(dbError)
	}

	log.Printf("Storing %s csv data in SQLite ...\n", dbName)
	for r := range csvData {
		address := csvData[r][0]
		room := csvData[r][1]
		apName := csvData[r][2]
		status := csvData[r][3]
		apType := csvData[r][4]
		apLoad := csvData[r][5]

		_, err := stmt.Exec(address, room, apName, status, apType, apLoad)
		if err != nil {
			panic(err)
		}
	}
	log.Printf("Finished data storing for %s.\n", dbName)
}

// creates a DB table (if not exists) to store an overview of all access points
func createTable(tableName string, db *sql.DB) {
	query := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS "%s" (
				"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"Address"	TEXT,
				"Room"		TEXT,
				"Name"		TEXT,
				"Status"	TEXT,
				"Type"		TEXT,
				"Load" 		TEXT
			);
		`, tableName)

	stmt, _ := db.Prepare(query)
	_, err := stmt.Exec()
	if err != nil {
		panic(err)
	}
}

// adds a new column of type 'TEXT' to table 'tableName'
func AddNewColumn(tableName, newCol string) {
	dbPath := fmt.Sprintf("%s%s.db", sqlitePath, tableName)
	db := initDB(dbPath)
	query := fmt.Sprintf(`
		ALTER TABLE %s ADD COLUMN %s TEXT
	`, tableName, newCol)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Panicf("Error: failed to prepare query! %q", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Panicf("Error: could not alter table! Failed to add new column: %s", err)
	}
}

func UpdateColumn(tableName, column , newValue, where string) {
	dbPath := fmt.Sprintf("%s%s.db", sqlitePath, tableName)
	db := initDB(dbPath)
	query := fmt.Sprintf(`UPDATE %s SET %s = ? WHERE %s%s`, tableName, column, column, where)
	stmt, dbError := db.Prepare(query)

	if dbError != nil {
		panic(dbError)
	}

	_, err := stmt.Exec(newValue)
	if err != nil {
		panic(err)
	}
}

func UpdateColumnName(tableName, currName, newName string) {
	dbPath := fmt.Sprintf("%s%s.db", sqlitePath, tableName)
	db := initDB(dbPath)
	query := fmt.Sprintf(`ALTER TABLE %s RENAME COLUMN %s TO %s;`, tableName, currName, newName)
	stmt, dbError := db.Prepare(query)

	if dbError != nil {
		panic(dbError)
	}

	_, err := stmt.Exec()
	if err != nil {
		panic(err)
	}
}
