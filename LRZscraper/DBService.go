package LRZscraper

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type AccessPointOverview struct {
	ID      string
	Address string
	Room    string
	Name    string
	Status  string
	Type    string
	Load    string
}

// initializes a local database instance, located in dbPath
// returns a pointer to the initialized database
func InitDB(dbPath string) *sql.DB {
	db, sqlError := sql.Open("sqlite3", dbPath)
	if sqlError != nil {
		panic(sqlError)
	}
	if db == nil {
		panic("db is nil")
	}
	return db
}

func CreateTableAccesspoints(db *sql.DB) {
	stmt, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "accesspoints" (
			"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"network"	TEXT,
			"current"	TEXT,
			"max"		TEXT,
			"min"		TEXT,
			"other" 	TEXT
		);
	`)
	_, err := stmt.Exec()
	if err != nil {
		panic(err)
	}
}

func ReadItem(db *sql.DB) []AccessPointOverview {
	query := `
		SELECT Address, Room, Load 
		FROM overview
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
		err2 := rows.Scan(&item.Address, &item.Room, &item.Load)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	return result
}

// for csv/graphData.csv
func storeDataInSQLite(dbPath string) {
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

	db := InitDB(dbPath)

	CreateTableAccesspoints(db)

	// fmt.Println(network, current, max, min, other)
	stmt, dbError := db.Prepare(`
		INSERT INTO accesspoints (network, current, max, min, other) values (?,?,?,?,?)
	`)

	if dbError != nil {
		panic(dbError)
	}

	fmt.Println("Storing data in SQLite ...")
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
	fmt.Println("Finished data storing")
}

// for csv/overview.csv
func storeOverviewInSQLite(dbPath string) {
	csvFile, err := os.Open("csv/overview.csv")
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

	db := InitDB(dbPath)
	CreateTableOverview(db)

	stmt, dbError := db.Prepare(`
		INSERT INTO overview (Address, Room, Name, Status, Type, Load) 
		values (?,?,?,?,?,?)
	`)

	if dbError != nil {
		panic(dbError)
	}

	fmt.Println("Storing data in SQLite ...")
	for r := range data {
		address := data[r][0]
		roomNr := data[r][1]
		apName := data[r][2]
		status := data[r][3]
		apType := data[r][4]
		apLoad := data[r][5]

		_, err := stmt.Exec(address, roomNr, apName, status, apType, apLoad)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Finished data storing")
}

// creates a DB table (if not exists) to store an overview of all access points
func CreateTableOverview(db *sql.DB) {
	stmt, _ := db.Prepare(`
			CREATE TABLE IF NOT EXISTS "overview" (
				"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				"Address"	TEXT,
				"Room"		TEXT,
				"Name"		TEXT,
				"Status"	TEXT,
				"Type"		TEXT,
				"Load" 		TEXT
			);
		`)
	_, err := stmt.Exec()
	if err != nil {
		panic(err)
	}
}
