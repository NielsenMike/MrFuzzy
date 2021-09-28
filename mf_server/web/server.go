package web

import (
	"container/list"
	"fmt"
	"html/template"
	"mf_server/data"
	"mf_server/sqlite"
	"mf_server/tarreader"
	"net/http"
)

var (
	ArchiveDirPtr *string = nil
)

func IndexHandler(w http.ResponseWriter, r *http.Request){
	tmpl := template.Must(template.ParseFiles("./web/html/index.html"))
	data := Data.DeviceInfo{DeviceName: "RaspberryPi - Wipro"} // TODO: Hardcoded Name
	data.Databases = Data.FindFilesByExtension(*ArchiveDirPtr, ".db")
	tmpl.Execute(w, data)
}


func DatabaseHandler(w http.ResponseWriter, r *http.Request){
	tmpl := template.Must(template.ParseFiles("./web/html/database.html"))
	dbAbsolutePath := r.URL.Query().Get("db")
	var dbPtr = sqlite.OpenDatabase(dbAbsolutePath)
	if dbPtr != nil{
		data := Data.DatabaseInfo{DatabaseName: dbAbsolutePath, DatabaseEntries: new([]Data.FileHashingDataSQL)}
		sqlite.SelectHashData(dbPtr, data.DatabaseEntries)
		tmpl.Execute(w, data)
		data.DatabaseEntries = nil
		dbPtr.Close()
	}
}

func WriteDataIntoDBHandler(w http.ResponseWriter, r *http.Request) {
	// open tar files
	for _, filename := range Data.FindFilesByExtension(*ArchiveDirPtr, ".tar") {
		tarReaderPtr := tarreader.Open(filename)
		data := list.New()
		err := tarreader.Read(tarReaderPtr, data) // read tar file and return list of file hashing data
		if err != nil {
			fmt.Println("Reading .tar file failed - ERROR: " + err.Error())
		}
		dbAbsolutePath := Data.ReplaceExt(filename, "db")
		var dbPtr = sqlite.OpenDatabase(dbAbsolutePath)
		if dbPtr != nil{
			sqlite.CreateTableHashedIfNotExists(dbPtr)
			for entry := data.Front(); entry != nil; entry = entry.Next() {
				name := entry.Value.(Data.FileHashingData).Name
				var hashSQLData = sqlite.SelectHashDataByName(dbPtr, name)
				if hashSQLData.Name == ""  {
					sqlite.InsertHashData(dbPtr, entry.Value.(Data.FileHashingData))
					continue
				}
				sqlite.UpdateHashData(dbPtr, entry.Value.(Data.FileHashingData))
			}
			//close database
			dbPtr.Close()
			fmt.Println("Finished All Chores Successfully")
		}
	}
}


