package web

import (
	"container/list"
	"fmt"
	"html/template"
	"mf_server/data"
	"mf_server/sqlite"
	"mf_server/tarreader"
	"net/http"
	"strconv"
)

var (
	ArchiveDirPtr *string = nil
	searchDataQuery = Data.FileHashingData{}
	databaseViewerEntriesNumber = 1000
)

func IndexHandler(w http.ResponseWriter, r *http.Request){
	tmpl := template.Must(template.ParseFiles("./web/html/index.html"))
	data := Data.ResourceInfo{
		Databases: Data.RetrieveFileInformation(*ArchiveDirPtr, ".db"),
		Archives: Data.RetrieveFileInformation(*ArchiveDirPtr, ".tar"),
		Backups: Data.RetrieveFileInformation(*ArchiveDirPtr, ".bak"),
	}
	tmpl.Execute(w, data)
}

func DatabaseHandler(w http.ResponseWriter, r *http.Request){
	var tmpl = template.Must(template.ParseFiles("./web/html/database.html"))
	var dbAbsolutePathGet = r.URL.Query().Get("db")
	var fromIndexGet = r.URL.Query().Get("from")
	var data = Data.DatabaseInfo{}
	var dbPtr = sqlite.OpenDatabase(dbAbsolutePathGet)
	if dbPtr != nil{
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			searchDataQuery = Data.ParseSearchString(r.FormValue("search"))
		}
		fromIndex, err := strconv.Atoi(fromIndexGet)
		if err != nil {
			fromIndex = 0
		}
		count := sqlite.SelectHashDataCountData(dbPtr, searchDataQuery)
		previousIndex, nextIndex := Data.GetPreviousNextIndex(fromIndex, databaseViewerEntriesNumber, count)
		data = Data.DatabaseInfo{DatabaseName: dbAbsolutePathGet, Entries: new([]Data.FileHashingDataSQL),
			Count: count, PreviousIndex: previousIndex, FromIndex: fromIndex, NextIndex: nextIndex}
		sqlite.SelectHashDataBySearch(dbPtr, searchDataQuery, fromIndex, databaseViewerEntriesNumber, data.Entries)
	}
	tmpl.Execute(w, data)
	data.Entries = nil
	dbPtr.Close()
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
		if !Data.FileExists(dbAbsolutePath) { Data.CreateFile(dbAbsolutePath) }
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


