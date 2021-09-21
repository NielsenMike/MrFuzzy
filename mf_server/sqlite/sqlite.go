package sqlite

import (
	"container/list"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"mf_server/data"
	"os"
	"strconv"
	"time"
)

/*
This function creates the sql table for the given database file
Defines the Database Scheme
*/
func createTable(database *sql.DB){

	//Creates the table i.e. the scheme
	createHashingTable := "CREATE TABLE IF NOT EXISTS hashed (" +
		"name NOT NULL PRIMARY KEY, " +
		"size TEXT, " +
		"init_sha256hash TEXT, " +
		"init_ssdeephash TEXT, " +
		"init_date TEXT, " +
		"cur_sha256hash TEXT," +
		"cur_ssdeephash TEXT," +
		"cur_date TEXT" +
		"missing BOOLEAN" +
		"percentChange INT);"

	//Prepare Statement
	fmt.Println("Making Table")
	statement, err := database.Prepare(createHashingTable)
	if err != nil {
		fmt.Println("Error making table: \n "+err.Error())
	}

	//Exec
	statement.Exec()
	fmt.Println("Made Table Successfully")
}

/*
This function writes to the table i.e. the database for the INIT SCENARIO
Filling the information with a prepared statement and a for loop
*/
func writeTableInit(database *sql.DB, data *list.List){

	//Writing the INIT data into the Table
	fmt.Println("Writing INIT Data")
	insertData := "INSERT INTO hashed(name, size, init_sha256Hash, init_ssdeepHash, init_date) VALUES (?,?,?,?,?)"
	statement, err := database.Prepare(insertData)
	if err != nil {
		fmt.Println("Error writing into database: \n" + err.Error())
	}

	//Date variable for init_date set
	currentTime := time.Now()
	var init_date = currentTime.Format("02-01-2006 15:04:05")

	//For looping through data and adding values from list
	for f := data.Front(); f != nil; f = f.Next() {
		name := f.Value.(Data.FileHashingData).Name
		size := f.Value.(Data.FileHashingData).Size
		sha256 := f.Value.(Data.FileHashingData).SHA256Hash
		ssdeep := f.Value.(Data.FileHashingData).SSDEEPHash
		statement.Exec(name, size, sha256, ssdeep, init_date)
	}
	fmt.Println("INIT Write Complete")
}

func WriteDatabase(filename string,data *list.List){

	//Creates a db file with table hashed
	fmt.Println("Creating database")
	file, err := os.Create("./hashing.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	_ = file.Close()
	fmt.Println("Created database")

	database, err := sql.Open("sqlite3", "./hashing.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	//CreateTable Database --> this could be an if statement so we don't always create a new table
	createTable(database)

	//WriteTable Database
	writeTableInit(database, data)

}

//Query Database for Hashing Data
func ReadDatabase(){
	database, err := sql.Open("sqlite3", "./hashing.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	//Reading
	fmt.Println("Reading")
	row, err := database.Query("SELECT * FROM hashed")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer row.Close()

	for row.Next(){
		var id int
		var name string
		var size string
		var sha256Hash string
		var ssdeepHash string
		row.Scan(&id, &name, &size, &sha256Hash, &ssdeepHash)
		fmt.Println(strconv.Itoa(id) + ": " + name + ": " + size + ": " + sha256Hash + ": " + ssdeepHash)
		fmt.Println("Read Complete")
	}
}
