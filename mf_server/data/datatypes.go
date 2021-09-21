package Data

import (
	"container/list"
	"fmt"
)


// file hashing data
type FileHashingData struct {
	Name string
	Size  string
	SHA256Hash string
	SSDEEPHash string
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