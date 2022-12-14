package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"mf_server/data"
	"os"
	"path/filepath"
)


func OpenDatabase(absolutePath string) *sql.DB{
	database, err := sql.Open("sqlite3", absolutePath)
	if err != nil {
		fmt.Println("Error opening db file: \n" + err.Error())
		return nil
	}
	return database
}

//createTable This function creates the sql table for the given database file
//Defines the Database Scheme/*
func CreateTableHashedIfNotExists(database *sql.DB){

	//Creates the table i.e. the scheme
	createHashingTable := "CREATE TABLE IF NOT EXISTS hashed (" +
		"name TEXT NOT NULL PRIMARY KEY, " +
		"size TEXT, " +
		"init_sha256hash TEXT, " +
		"init_ssdeephash TEXT, " +
		"init_date TEXT, " +
		"cur_sha256hash TEXT," +
		"cur_ssdeephash TEXT," +
		"cur_date TEXT," +
		"percentChange INT);"
	statement, err := database.Prepare(createHashingTable)
	if err != nil {
		fmt.Println("Error creating table: \n "+err.Error())
	}
	statement.Exec()
}

//createTable stats
func CreateTableStatsIfNotExists(database *sql.DB){
	createStatsTable := "CREATE TABLE IF NOT EXISTS stats (" +
		"lastUpdate TEXT);"
	statement, err := database.Prepare(createStatsTable)
	if err != nil {
		fmt.Println("Error creating table: \n "+err.Error())
	}
	statement.Exec()
}

func SelectHashDataByName(database *sql.DB, name string) Data.FileHashingDataSQL{
	var hashedSQLData = Data.FileHashingDataSQL{}
	err := database.QueryRow("SELECT * FROM hashed WHERE name = ?", name).Scan(&hashedSQLData.Name,
		&hashedSQLData.Size, &hashedSQLData.InitSha256hash, &hashedSQLData.InitSsdeephash, &hashedSQLData.InitDate,
		&hashedSQLData.CurSha256hash, &hashedSQLData.CurSsdeephash, &hashedSQLData.CurDate, &hashedSQLData.PercentChange)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		fmt.Println("Error in select hash data \n" + err.Error())
	}
	return hashedSQLData
}

func SelectDatabaseStatsUpdateTime(database *sql.DB) string{
	var updateDate = ""
	err := database.QueryRow("SELECT max(lastUpdate) from stats").Scan(&updateDate)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		fmt.Println("Error in select hash data \n" + err.Error())
	}
	return updateDate
}

func SelectDatabaseStatsInitTime(database *sql.DB) string{
	var initDate = ""
	err := database.QueryRow("SELECT MIN(lastUpdate) from stats").Scan(&initDate)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		fmt.Println("Error in select hash data \n" + err.Error())
	}
	return initDate
}

func SelectHashDataCountData(database *sql.DB, searchAttributes Data.FileHashingData) int{
	var count = 0
	countQuery := "SELECT COUNT(*) FROM hashed WHERE " +
		"name like ? " +
		"AND (init_sha256hash LIKE ? " +
		"OR cur_sha256hash LIKE ? )" +
		"AND (init_ssdeephash LIKE ? " +
		"OR cur_ssdeephash LIKE ? )"
	countStatement, err := database.Prepare(countQuery)
	if err != nil{
		fmt.Println("Error select statements: \n" + err.Error())
		return count
	}
	sqlRow := countStatement.QueryRow("%" + searchAttributes.Name + "%",
		"%" + searchAttributes.SHA256Hash + "%",
		"%" + searchAttributes.SHA256Hash + "%",
		"%" + searchAttributes.SSDEEPHash + "%",
		"%" + searchAttributes.SSDEEPHash + "%")
	sqlRow.Scan(&count)
	return count
}

func SelectHashDataBySearch(database *sql.DB, searchAttributes Data.FileHashingData, fromIndex, maxRows int,
	outHashingData *[]Data.FileHashingDataSQL) {
	selectQuery := "SELECT * FROM hashed WHERE " +
		"name like ? " +
		"AND (init_sha256hash LIKE ?" +
		"OR cur_sha256hash LIKE ? )" +
		"AND (init_ssdeephash LIKE ? " +
		"OR cur_ssdeephash LIKE ? )" +
		" ORDER BY percentChange ASC, cur_date ASC, init_date DESC LIMIT ?,?"
	selectStatement, err := database.Prepare(selectQuery)
	if err != nil{
		fmt.Println("Error select statements: \n" + err.Error())
		return
	}
	sqlRows, err := selectStatement.Query("%" + searchAttributes.Name + "%",
		"%" + searchAttributes.SHA256Hash + "%",
		"%" + searchAttributes.SHA256Hash + "%",
		"%" + searchAttributes.SSDEEPHash + "%",
		"%" + searchAttributes.SSDEEPHash + "%",
		fromIndex,
		maxRows)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Error in select hash data \n" + err.Error())
	}else{
		for sqlRows.Next() {
			var selRow Data.FileHashingDataSQL
			if err := sqlRows.Scan(&selRow.Name,
				&selRow.Size, &selRow.InitSha256hash, &selRow.InitSsdeephash, &selRow.InitDate,
				&selRow.CurSha256hash, &selRow.CurSsdeephash, &selRow.CurDate, &selRow.PercentChange);
				err != nil {
				return
			}
			*outHashingData = append(*outHashingData, selRow)
		}
	}
}

func InsertStatsData(database *sql.DB, currentDateTime string) {
	insertData := "INSERT INTO stats (lastUpdate) values(?)"
	statement, err := database.Prepare(insertData)
	if err != nil {
		fmt.Println("Error insert into database: \n" + err.Error())
		return
	}
	statement.Exec(currentDateTime)
}

func InsertHashData(database *sql.DB, data Data.FileHashingData, initDate string) {
	insertData := "INSERT INTO hashed (name, size, init_sha256Hash, init_ssdeepHash, init_date, cur_sha256hash, " +
		"cur_ssdeephash, cur_date, percentchange) VALUES (?,?,?,?,?,?,?,?,?)"
	statement, err := database.Prepare(insertData)
	if err != nil {
		fmt.Println("Error insert into database: \n" + err.Error())
		return
	}
	statement.Exec(data.Name, data.Size, data.SHA256Hash, data.SSDEEPHash, initDate, "null","null","null", 0)
}


//writeTableCur This function writes to the table i.e. the database for the UPDATE SCENARIO
//Filling the information with a prepared statement and a for loop/*
func UpdateHashData(database *sql.DB, data Data.FileHashingData, curDate string){
	var sqlHashData = SelectHashDataByName(database, data.Name)
	updateStatement := "UPDATE hashed SET cur_sha256hash = $1, cur_ssdeephash = $2, cur_date =  $3, percentChange = $4 WHERE name = $5"
	var percentageChanged = Data.CalculateMatchScore(sqlHashData.InitSsdeephash,
		data.SSDEEPHash, sqlHashData.InitSha256hash, data.SHA256Hash, sqlHashData.PercentChange)
	_, err := database.Exec(updateStatement, data.SHA256Hash, data.SSDEEPHash, curDate, percentageChanged, data.Name)
	if err != nil {
		fmt.Println("Error writing into database: \n" + err.Error())
	}
}

// Cleanup Serves as the cleanup, removes the remaining TAR files from the directory
//Leaving the db file/*
func Cleanup(finalpath string){

	files, err := ioutil.ReadDir(finalpath)
	if err != nil {
		fmt.Println("Error reading given Directory: \n" + err.Error())
	}
	for _, file := range files {
		var foundtar, _ = filepath.Match("*.tar", file.Name())
		if foundtar {
			var fullfilename string = finalpath +"/"+ file.Name()
			err := os.Remove(fullfilename)
			if err != nil {
				fmt.Println("Removing file error: \n"+ err.Error())
			}
		}
	}

}


