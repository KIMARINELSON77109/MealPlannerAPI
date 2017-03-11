package main

import (
	"database/sql"
	"log"
)

type Uom struct {
	Uom 	string `json:"uom"`
	UomId 	int	   `json:"uomId"`
}

func SaveUom (dbh sql.DB, uom Uom) (error) {
	//TODO:  check if the uom is already in the db
	stmt, err := dbh.Prepare("insert into uom(uom_id, uom) values(?, ?) on duplicate key update set uom_id=?, uom=?")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(uom.UomId, uom.Uom, uom.UomId, uom.Uom)
	if err != nil {
		log.Printf("Error saving UOM: %v\n", err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting last id for UOM: %v\n", err)
	}
	log.Printf("Added uom id of: %d", lastId)
	stmt.Close()
	return err
}

func GetUom (dbh sql.DB, uom string) (Uom, error) {
	// returns the Uom matched on the uom string 
	stmt, err := dbh.Prepare("select uom, uom_id from uom where uom = ?")
	if err != nil {
		log.Fatal(err)
	}
	var myUom = new(Uom)
	err = stmt.QueryRow(uom).Scan(&myUom.Uom, &myUom.UomId)
	if err != nil {
		log.Printf("Error getting Uom: %v\n", err)
	}
	stmt.Close()
	return *myUom, err
}
