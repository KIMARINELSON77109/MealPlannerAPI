package main
// the kick off for the api
// Fred T. Dunaway
// August 25, 2016

import (
	"log"
    "net/http"
)

const mysqlTimeFormat = "2006-01-02 15:04:05"
//var dbh *sql.DB		// our global database handler pointer.

func main() {

    router := NewRouter()
	certFile := "./server.pm"
    keyFile := "./server.key"
    
//    log.Fatal(http.ListenAndServe(":8080", router))
	log.Fatal(http.ListenAndServeTLS(":3000", certFile, keyFile, router))
}
