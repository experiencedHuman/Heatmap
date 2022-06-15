package DBService

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"fmt"
)

func SetupHistoryTable() {
	tableName := "history"
	CreateHistoryTable(tableName)
	apList := RetrieveAPsOfTUM(true)
	PopulateHistoryTable(apList)
}

func HistoryCSVtoSQLite() {
	file, err := os.Open("./data/csv/histories2.csv")
	if err != nil {
		log.Println("Could not open file!", err)
	}

	csvReader := csv.NewReader(file)
	csvReader.Comma = ','

	for {
		record, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("couldn't read record")
			continue
		}
		
		apName := record[0]
		max := record[1]
		min := record[2]

		where := fmt.Sprintf("Name='%s'", apName)
		UpdateColumn("apstat","Max",max, where)
		UpdateColumn("apstat","Min",min, where)
	}
}

// JOIN
func JoinMaxMin() {
	// query := `
	// 	UPDATE history
	// 	SET history.Max = apstat.Max, history.Min = apstat.Min
	// 	FROM history INNER JOIN apstat ON history.AP_Name = apstat.Name
	// `
	
	// max
	// query := `
	// 	UPDATE history
	// 	SET Max = (SELECT Max 
	// 	FROM apstat WHERE history.AP_Name = Name)
	// `

	// min
	query := `
		UPDATE history
		SET Min = (SELECT Min
		FROM apstat WHERE history.AP_Name = Name)
	`
	runQuery(query)
}

// READ AP HISTORY with max min values
func ReadHistoryWithMaxMin(apName string, day int, hour int) {

}