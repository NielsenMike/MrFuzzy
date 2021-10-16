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

type FileInformation struct {
	AbsolutePath string
	Timestamp string
	Size string
}


// page data for Resources
type ResourceInfo struct {
	Databases []FileInformation
	Archives []FileInformation
	Backups []FileInformation
}

// page data for database entries
type DatabaseInfo struct {
	DatabaseName string
	Entries *[]FileHashingDataSQL
	Count int
	PreviousIndex int
	FromIndex int
	NextIndex int
	InitializedDate string
	LastUpdateDate string
}

func GetPreviousNextIndex(fromIndex, size, count int) (int,int){
	var nextIndex = fromIndex + size
	var previousIndex = fromIndex - size
	if nextIndex > count{
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
	ssdeepHash, err := ssdeep.FuzzyBytes(*data)
	if err != nil {
		ssdeepHash = "" // the file is to small for ssdeep output error
	}
	sha256hash := sha256.New()
	sha256hash.Write(*data) // write bytes into sha object

	fileHashingData.SSDEEPHash = ssdeepHash
	fileHashingData.SHA256Hash = hex.EncodeToString(sha256hash.Sum(nil))
}

func CalculateMatchScore(ssdeep1, ssdeep2, sha1, sha2 string, currentScore int) int {
	score := currentScore
	if len(ssdeep1) > 0 && len(ssdeep2) > 0 {
		score, _ = ssdeep.Distance(ssdeep1, ssdeep2)
	} else {
		if strings.Compare(sha1, sha2) == 0{
			score = 100
		}else { score = 0 }
	}
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
		case "sha":
			fileHashingData.SHA256Hash = pair[1]
		case "ssdeep":
			fileHashingData.SSDEEPHash = pair[1]
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