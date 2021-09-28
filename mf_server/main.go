package main

import (
	"flag"
	"log"
	"mf_server/web"
	"net/http"
)



func main() {
	archiveDirPtr := flag.String("archdir", "", "a string")
	flag.Parse()
	if *archiveDirPtr != ""{
		web.ArchiveDirPtr = archiveDirPtr
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
		http.HandleFunc("/", web.IndexHandler)
		http.HandleFunc("/database", web.DatabaseHandler)
		http.HandleFunc("/writeDataIntoDB", web.WriteDataIntoDBHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}




