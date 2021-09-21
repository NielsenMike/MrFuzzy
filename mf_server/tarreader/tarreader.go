package tarreader

import (
	"archive/tar"
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"github.com/glaslos/ssdeep"
	"io"
	"log"
	"mf_server/data"
	"strconv"
)

// read each file from tar
func Read(tr *tar.Reader, ls *list.List) error {
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

func readFiles(header *tar.Header, data []byte) Hashing.FileHashingData {
	// if it's a file create it
	fileData := Hashing.FileHashingData{}
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
		fileData = Hashing.FileHashingData{
			Name: header.Name,
			Size: strconv.FormatInt(header.Size, 10),
			SHA256Hash: hex.EncodeToString(sha256hash.Sum(nil)),
			SSDEEPHash: ssdeepHash,
		}
	}
	// return file hashing data element
	return fileData
}
