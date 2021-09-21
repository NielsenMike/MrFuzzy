package main

import (
	"container/list"
	"flag"
	"fmt"
	"mf_server/data"
	"mf_server/sqlite"
	"mf_server/tarreader"
)

func main() {

	// parse arguments
	tarDirPtr := flag.String("tardir", "", "a string")
	dbDirPtr := flag.String("dbdir", "", "a string")
	flag.Parse()

	// check if both arguments are set
	if *tarDirPtr != "" && *dbDirPtr != "" {
		// open tar files
		for _, filename := range Data.Find(*tarDirPtr, ".tar") {
			tarReaderPtr := tarreader.Open(filename)
			data := list.New()
			err := tarreader.Read(tarReaderPtr, data) // read tar file and return list of file hashing data
			if err != nil {
				fmt.Println("Reading .tar file failed - ERROR: " + err.Error())
			}
			dbFilename := Data.ReplaceExt(filename, "db")
			//Store Data in SQLite Format
			sqlite.WriteDatabase(dbFilename, data)
		}
	}
}



