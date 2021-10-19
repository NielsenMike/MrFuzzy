package web

import (
	"container/list"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mf_server/data"
	"mf_server/sqlite"
	"mf_server/tarreader"
	"net/http"
	"os"
	"strconv"
	"time"
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
		initializedDate := sqlite.SelectDatabaseStatsInitTime(dbPtr)
		lastUpdatedDate := sqlite.SelectDatabaseStatsUpdateTime(dbPtr)
		count := sqlite.SelectHashDataCountData(dbPtr, searchDataQuery)
		previousIndex, nextIndex := Data.GetPreviousNextIndex(fromIndex, databaseViewerEntriesNumber, count)
		data = Data.DatabaseInfo{DatabaseName: dbAbsolutePathGet,
			Entries: new([]Data.FileHashingDataSQL),
			Count: count,
			PreviousIndex: previousIndex,
			FromIndex: fromIndex,
			NextIndex: nextIndex,
			InitializedDate: initializedDate,
			LastUpdateDate: lastUpdatedDate,
		}
		sqlite.SelectHashDataBySearch(dbPtr, searchDataQuery, fromIndex, databaseViewerEntriesNumber, data.Entries)
	}
	tmpl.Execute(w, data)
	data.Entries = nil
	dbPtr.Close()
}

func WriteDataIntoDBHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	var currentDateTime = currentTime.Format("02-01-2006 15:04:05")
	for _, tarAbsolutePath := range Data.FindFilesByExtension(*ArchiveDirPtr, ".tar") {
		fmt.Printf("Reading Tar-file: %s \n", tarAbsolutePath)
		data := list.New()
		err := readDataFromTar(tarAbsolutePath, data)
		if err == nil{
			dbAbsolutePath := Data.ReplaceExt(tarAbsolutePath, "db")
			if !Data.FileExists(dbAbsolutePath) {
				Data.CreateFile(dbAbsolutePath)
			}
			err = updateDatabaseChanges(dbAbsolutePath, data, currentDateTime)
			if err == nil{
				fmt.Printf("Filled data to database: %s \n", dbAbsolutePath)
				err = updateDatabaseStats(dbAbsolutePath, currentDateTime)
				if err == nil{
					fmt.Printf("Updated stats to database: %s \n", dbAbsolutePath)
					err = createBackupFromTar(tarAbsolutePath)
					if err == nil{
						fmt.Printf("Created backup file: %s \n", tarAbsolutePath)
					}else{ fmt.Println("Failed to create backup " + err.Error()) }
				}else{
					fmt.Printf("Failed to update stats for database: %s \n", dbAbsolutePath)
				}
			} else { fmt.Println("Failed to fill data to database " + err.Error()) }
		} else{ fmt.Println("Failed read data from tar " + err.Error()) }
	}
}

func createBackupFromTar(absolutePath string) error{
	// copy tar create backup
	backupAbsolutePath := Data.ReplaceExt(absolutePath, strconv.FormatInt(time.Now().Unix(), 10) + ".tar.bak")
	source, err := os.Open(absolutePath)
	if err != nil {
		return errors.New("Failed to open file.")
	}
	destination, err := os.Create(backupAbsolutePath)
	if err != nil {
		return errors.New("Failed to create file.")
	}
	_, err = io.Copy(destination, source)
	if err != nil{
		return errors.New("Failed to copy file.")
	}
	destination.Close()
	source.Close()
	return nil
}

func readDataFromTar(absolutePath string, dataOut *list.List) error{
	tarReaderPtr := tarreader.Open(absolutePath)
	err := tarreader.Read(tarReaderPtr, dataOut) // read tar file and return list of file hashing data
	if err != nil {
		return err
	}
	return nil
}

func updateDatabaseChanges(absolutePath string, inData* list.List, currentDateTime string) error{
	var dbPtr = sqlite.OpenDatabase(absolutePath)
	if dbPtr != nil {
		sqlite.CreateTableHashedIfNotExists(dbPtr)
		for entry := inData.Front(); entry != nil; entry = entry.Next() {
			name := entry.Value.(Data.FileHashingData).Name
			var hashSQLData = sqlite.SelectHashDataByName(dbPtr, name)
			if hashSQLData.Name == "" {
				sqlite.InsertHashData(dbPtr, entry.Value.(Data.FileHashingData), currentDateTime)
				continue
			}
			sqlite.UpdateHashData(dbPtr, entry.Value.(Data.FileHashingData), currentDateTime)
		}
		dbPtr.Close()
		return nil
	}
	return errors.New("failed to set data into database.")
}

func updateDatabaseStats(absolutePath, currentDateTime string) error{
	var dbPtr = sqlite.OpenDatabase(absolutePath)
	if dbPtr != nil {
		sqlite.CreateTableStatsIfNotExists(dbPtr)
		sqlite.InsertStatsData(dbPtr, currentDateTime)
		dbPtr.Close()
		return nil
	}
	return errors.New("failed to set stats into database.")
}


