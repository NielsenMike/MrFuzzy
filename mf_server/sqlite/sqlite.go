package sqlite

import (
	"container/list"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"mf_server/data"
	"os"
	"path/filepath"
	"strings"
	"time"
)


//createTable This function creates the sql table for the given database file
//Defines the Database Scheme/*
func createTable(database *sql.DB){

	//Creates the table i.e. the scheme
	createHashingTable := "CREATE TABLE hashed (" +
		"name NOT NULL PRIMARY KEY, " +
		"size TEXT, " +
		"init_sha256hash TEXT, " +
		"init_ssdeephash TEXT, " +
		"init_date TEXT, " +
		"cur_sha256hash TEXT," +
		"cur_ssdeephash TEXT," +
		"cur_date TEXT," +
		"missing TEXT," +
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


//writeTableInit This function writes to the table i.e. the database for the INIT SCENARIO
//Filling the information with a prepared statement and a for loop /*
func writeTableInit(database *sql.DB, data *list.List){

	//Writing the INIT data into the Table
	fmt.Println("Writing INIT Data")
	insertData := "INSERT INTO hashed (name, size, init_sha256Hash, init_ssdeepHash, init_date, cur_sha256hash, cur_ssdeephash, cur_date, missing, percentchange) VALUES (?,?,?,?,?,?,?,?,?,?)"
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
		statement.Exec(name, size, sha256, ssdeep, init_date, "null","null","null", "false", 100)
	}
	fmt.Println("INIT Write Complete")
}


//writeTableCur This function writes to the table i.e. the database for the UPDATE SCENARIO
//Filling the information with a prepared statement and a for loop/*
func writeTableCur(database *sql.DB, data *list.List){

	//Date variable for init_date set
	currentTime := time.Now()
	var cur_date = currentTime.Format("02-01-2006 15:04:05")

	//Writing the CUR data into the Table
	fmt.Println("Writing UPDATE Data")
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
		_, err := database.Exec(updateStatement, sha256, ssdeep, cur_date, "false", 100, name1)
		if err != nil {
			//ENTRY DOES NOT EXIST --> Newly added file from Update by threat actor or system
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

// WriteDatabase This function serves to be the "main" function of this package
//The logical process is decided whether a db file exists and the UPDATE SCENARIO is used or
//if the INIT SCENARIO is used, when no db file exists.
//DB file is created
//Cleanup is also performed/*
func WriteDatabase(filename string,data *list.List){

	//init or update boolean
	var doinit = true

	//Get Path from filename
	var pathname  = strings.Split(filename, "/")
	var dbfile string = pathname[len(pathname)-1]
	pathname = pathname[:len(pathname)-1]
	var finalpath = strings.Join(pathname, "/")

	//Start logical decision: init or update
	//List files within given filepath
	files, err := ioutil.ReadDir(finalpath)
	if err != nil {
		fmt.Println("Error reading given Directory: \n" + err.Error())
	}
	for _, file := range files {
		var founddb, _ = filepath.Match(dbfile, file.Name())
		if founddb {
			doinit = false
			break
		}
	}
	//Decision: init or update
	if doinit {
		//Did not find db --> init process
		//Create db file
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error making db file: \n"+err.Error())
		}
		_ = file.Close()

		database, err := sql.Open("sqlite3", filename)
		if err != nil {
			fmt.Println("Error opening db file: \n" + err.Error())
		}
		defer database.Close()

		//CreateTable Database
		createTable(database)

		//WriteTable Database
		writeTableInit(database, data)

		//close database
		database.Close()

	} else {
		database, err := sql.Open("sqlite3", filename)
		if err != nil {
			fmt.Println("Error opening db file: \n" + err.Error())
		}
		defer database.Close()

		//updateTable
		writeTableCur(database, data)

		//close database
		database.Close()
	}

	fmt.Println("Finished All Chores Successfully")
}


