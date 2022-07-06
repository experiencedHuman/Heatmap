package DBService

import (
	"fmt"
	"log"
)

func GetAllNames() []string {
	query := `
		SELECT Name
		FROM apstat
		WHERE Address LIKE '%TUM%'
			AND Lat!='lat'
			AND Long!='long'
	`
	rows, err := DB.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)

		if err != nil {
			panic(err)
		}
		names = append(names, name)
	}
	return names
}

func GetAccessPointByName(name string) *AccessPoint {
	stmt := fmt.Sprintf(`
		SELECT Name, Lat, Long, Load
		FROM apstat
		WHERE Name='%s'
	`, name)

	row := DB.QueryRow(stmt)
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

// Stores csv data in SQLite table
// under path parameter 'dbPath'.
func StoreApstatInDB(csvPath, dbPath, tableName string) {
	csvData := readFromCSV(csvPath)
	db := InitDB(dbPath)

	createTable := fmt.Sprintf(
		`
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
