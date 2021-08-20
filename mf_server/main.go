package main

import (
	"archive/tar"
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/glaslos/ssdeep"
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
	// output of all file hashing data of tar file
	for f := data.Front(); f != nil; f = f.Next() {
		fmt.Println(f.Value.(fileHashingData).Name +
			" " + f.Value.(fileHashingData).Size +
			" " + f.Value.(fileHashingData).sha256Hash +
			" " + f.Value.(fileHashingData).ssdeepHash,
		)
	}
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

