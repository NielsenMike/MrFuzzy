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
)

//create Table function
func createTable(database *sql.DB){
	createHashingTable := "CREATE TABLE IF NOT EXISTS hashed (" +
		"id INTEGER PRIMARY KEY, " +
		"name TEXT, " +
		"size TEXT, " +
		"sha256Hash TEXT, " +
		"ssdeepHash TEXT);"
	fmt.Println("Making Table")
	statement, err := database.Prepare(createHashingTable)
	if err != nil {
		fmt.Println(err.Error())
	}
	statement.Exec()
	fmt.Println("Made Table")
}

func writeTable(database *sql.DB, data *list.List){

	fmt.Println("Writing to Database...")
	insertData := "INSERT INTO hashed(name, size, sha256Hash, ssdeepHash) VALUES (?,?,?,?)"
	statement, err := database.Prepare(insertData)
	if err != nil {
		fmt.Println(err.Error())
	}
	//For looping through data and adding values
	for f := data.Front(); f != nil; f = f.Next() {
		name := f.Value.(Hashing.FileHashingData).Name
		size := f.Value.(Hashing.FileHashingData).Size
		sha256 := f.Value.(Hashing.FileHashingData).SHA256Hash
		ssdeep := f.Value.(Hashing.FileHashingData).SSDEEPHash
		statement.Exec(name, size, sha256, ssdeep)
	}
	fmt.Println("Write Complete")
}

//Write Hashing Data into SQLLite Database --> READ TODOS
func WriteDatabase(data *list.List){

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
	writeTable(database, data)

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
