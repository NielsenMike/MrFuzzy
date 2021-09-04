package main

import (
	"archive/tar"
	"container/list"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/glaslos/ssdeep"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	// current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	// define a tar file to analyze, should be in the same folder as this file
	file, err := os.Open(workingDir + "/archive.tar")
	tr := tar.NewReader(file)
	data := list.New()
	err = readTar(tr, data) // read tar file and return list of file hashing data
	if err != nil{
		fmt.Println(err) // something went wrong
	}

	/*
	// output of all file hashing data of tar file
	for f := data.Front(); f != nil; f = f.Next() {
		fmt.Println(f.Value.(fileHashingData).Name +
			" " + f.Value.(fileHashingData).Size +
			" " + f.Value.(fileHashingData).sha256Hash +
			" " + f.Value.(fileHashingData).ssdeepHash,
		)
	}
	*/

	//Store Data in SQLite Format
	writeDatabase(data)

	//Read Database
	readDatabase()

}

// file hashing data
type fileHashingData struct {
	Name string
	Size  string
	sha256Hash string
	ssdeepHash string

}

// read each file from tar
func readTar(tr *tar.Reader, ls *list.List) error {
	for {
		header, err := tr.Next()
		switch {
		// if no more files are found return
		case err == io.EOF:
			return nil
		// return any other error
		case err != nil:
			return err
		// if the header is nil, just skip it
		case header == nil:
			continue
		}
		data, err := io.ReadAll(tr)
		if err != nil {
			log.Fatal(err) // something went wrong with the file
		}
		fileData := readFiles(header, data)

		// if it is not a file don't add element into list
		if len(fileData.Name) > 0{
			ls.PushBack(fileData)
		}
	}
}

func readFiles(header *tar.Header, data []byte) fileHashingData{
	// if it's a file create it
	fileData := fileHashingData{}
	// check if it is a file
	if header.Typeflag == tar.TypeReg {
		// create fuzzy hash with ssdeep
		ssdeepHash, fError := ssdeep.FuzzyBytes(data)
		if fError != nil {
			ssdeepHash = fError.Error() // the file is to small for ssdeep output error
		}
		sha256hash := sha256.New()
		sha256hash.Write(data) // write bytes into sha object

		// create file hashing data struct
		fileData = fileHashingData{
			Name: header.Name,
			Size: strconv.FormatInt(header.Size, 10),
			sha256Hash: hex.EncodeToString(sha256hash.Sum(nil)),
			ssdeepHash: ssdeepHash,
		}
	}
	// return file hashing data element
	return fileData
}

//Write Hashing Data into SQLLite Database --> READ TODOS
func writeDatabase(data *list.List){

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
		statement.Exec(f.Value.(fileHashingData).Name, f.Value.(fileHashingData).Size,
			f.Value.(fileHashingData).sha256Hash, f.Value.(fileHashingData).ssdeepHash)
	}
	fmt.Println("Write Complete")
}

//Query Database for Hashing Data
func readDatabase(){
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

