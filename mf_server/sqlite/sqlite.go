package sqlite

import (
	"container/list"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"mf_server/data"
	"os"
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
	insertData := "INSERT INTO hashed(name, size, init_sha256Hash, init_ssdeepHash, init_date, cur_sha256hash, cur_ssdeephash, cur_date, missing, percentchange) VALUES (?,?,?,?,?,?,?,?,?,?)"
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
		statement.Exec(name, size, sha256, ssdeep, init_date, "null","null","null", false, 100)
	}
	fmt.Println("INIT Write Complete")
}

/*
This function writes to the table i.e. the database for the UPDATE SCENARIO
Filling the information with a prepared statement and a for loop
*/
func writeTableCur(database *sql.DB, data *list.List){

	//Date variable for init_date set
	currentTime := time.Now()
	var cur_date = currentTime.Format("02-01-2006 15:04:05")

	//Writing the CUR data into the Table
	fmt.Println("Writing INIT Data")
	updateStatement := "UPDATE hashed SET cur_sha256hash = $1, cur_ssdeephash = $2, cur_date =  $3, missing = $4, percentChange = $5 WHERE name = $6"


	//For looping through data and adding values from list - UPDATE LIST
	for f := data.Front(); f != nil; f = f.Next() {
		name1 := f.Value.(Data.FileHashingData).Name
		size := f.Value.(Data.FileHashingData).Size
		sha256 := f.Value.(Data.FileHashingData).SHA256Hash
		ssdeep := f.Value.(Data.FileHashingData).SSDEEPHash

		//Call to see if entry exists
		// YES --> Continue
		// NO --> Add entry as new
		//Missing and Percent Flags updated here --> TODO
		_, err := database.Exec(updateStatement, sha256, ssdeep, cur_date, false, 100, name1)
		if err != nil {
			//ENTRY DOES NOT EXIST --> Newly added file from Update
			if err == sql.ErrNoRows {

				//Inserting New file
				insertData := "INSERT INTO hashed(name, size, init_sha256Hash, init_ssdeepHash, init_date, cur_sha256hash, cur_ssdeephash, cur_date, missing, percentchange) VALUES (?,?,?,?,?,?,?,?,?,?)"
				statement, err := database.Prepare(insertData)
				if err != nil {
					fmt.Println("Error writing into database: \n" + err.Error())
				}
				statement.Exec(name1, size, sha256, ssdeep, cur_date, sha256, ssdeep, cur_date, false, 100)

			}
			fmt.Println("Error writing into database: \n" + err.Error())
		}
	}
	fmt.Println("CUR Write Complete")
}


/*
This function serves to be the "main" function of this package
The logcial process is decided as to whether a db file exists and the UPDATE SCENARIO is used or
if the INIT SCENARIO is used, when no db file exists.
DB file is created
Cleanup is also performed
*/
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


