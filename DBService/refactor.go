package DBService

func SetupHistoryTable() {
	tableName := "history"
	CreateHistoryTable(tableName)
	apList := RetrieveAPsOfTUM(true)
	PopulateHistoryTable(apList)
}