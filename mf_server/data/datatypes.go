package Data

import (
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/glaslos/ssdeep"
	"strings"
)


// file hashing data
type FileHashingData struct {
	Name string
	Size  string
	SHA256Hash string
	SSDEEPHash string
}

// file hashing data
type FileHashingDataSQL struct {
	Name           string
	Size           string
	InitSha256hash string
	InitSsdeephash string
	InitDate       string
	CurSha256hash  string
	CurSsdeephash  string
	CurDate        string
	PercentChange  int
}

// page data for device info
type DeviceInfo struct {
	DeviceName string
	Databases []string
}

// page data for database entries
type DatabaseInfo struct {
	DatabaseName string
	Entries *[]FileHashingDataSQL
	Count int
	PreviousIndex int
	FromIndex int
	NextIndex int
}

func GetPreviousNextIndex(fromIndex, size, count int) (int,int){
	var nextIndex = fromIndex + size
	var previousIndex = fromIndex - size
	if nextIndex >= count{
		nextIndex = count
	}
	if 0 >= previousIndex {
		previousIndex = 0
	}
	return previousIndex, nextIndex
}

// set SHA256 & SSDEEP hash values
func SetHashValues(fileHashingData *FileHashingData, data *[]byte){
	// create fuzzy hash with ssdeep
	ssdeepHash, fError := ssdeep.FuzzyBytes(*data)
	if fError != nil {
		ssdeepHash = "-" // the file is to small for ssdeep output error
	}
	sha256hash := sha256.New()
	sha256hash.Write(*data) // write bytes into sha object

	fileHashingData.SSDEEPHash = ssdeepHash
	fileHashingData.SHA256Hash = hex.EncodeToString(sha256hash.Sum(nil))
}

func CalculateSSDEEPScore(hash1 string, hash2 string, currentScore int) int {
	score := currentScore
	score, _ = ssdeep.Distance(hash1, hash2)
	return score
}


func ParseSearchString(searchString string) FileHashingData{
	fileHashingData := FileHashingData{}
	query := strings.Split(searchString, ";")
	for _, attribute := range query {
		pair := strings.Split(attribute, "=")
		switch pair[0] {
		case "name":
			fileHashingData.Name = pair[1]
		}
	}
	return fileHashingData
}


// print data
func PrintData(data *list.List)  {
		for f := data.Front(); f != nil; f = f.Next() {
			name := f.Value.(FileHashingData).Name
			size := f.Value.(FileHashingData).Size
			sha256 := f.Value.(FileHashingData).SHA256Hash
			ssdeep := f.Value.(FileHashingData).SSDEEPHash
			fmt.Println(name +
				" " + size +
				" " + sha256 +
				" " + ssdeep,
			)
	}
}