package tarreader

import (
	"archive/tar"
	"container/list"
	"io"
	"log"
	"mf_server/data"
	"os"
	"strconv"
)

func Open(filepath string) *tar.Reader{
	fd, err := os.Open(filepath)
	if err != nil{
		panic("Open .tar file failed on path:" +filepath)
	}
	tarReader := tar.NewReader(fd)
	return tarReader
}


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

func readFiles(header *tar.Header, data []byte) Data.FileHashingData {
	// if it's a file create it
	fileData := Data.FileHashingData{}
	// check if it is a file
	if header.Typeflag == tar.TypeReg {
		fileData.Name = header.Name
		fileData.Size = strconv.FormatInt(header.Size, 10)
		Data.SetHashValues(&fileData, &data)
	}
	// return file hashing data element
	return fileData
}
