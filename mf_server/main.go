package main

import (
	"archive/tar"
	"container/list"
	"fmt"
	"mf_server/tarreader"
	"os"
)

func main() {
	programName := os.Args[0]
	println(programName)


	// current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	// define a tar file to analyze, should be in the same folder as this file
	file, err := os.Open(workingDir + "/archive.tar")
	tr := tar.NewReader(file)
	data := list.New()
	err = tarreader.Read(tr, data) // read tar file and return list of file hashing data
	if err != nil{
		fmt.Println(err) // something went wrong
	}

	//Store Data in SQLite Format
	//sqlite.WriteDatabase(data)

	//Read Database
	//sqlite.ReadDatabase()
}


