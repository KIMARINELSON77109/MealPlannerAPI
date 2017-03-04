package main

import (
	"log"
	"database/sql"
//	"strconv"
)

type Owner struct {
	Name 	string `json:"name"`
	Email	string `json:"email"`
}

func GetOwner (dbh sql.DB, ownerEmail string) (Owner, error){	
	stmt, err := dbh.Prepare("select owner_name, owner_email from owner where owner_email=?")
	if err != nil {
		log.Fatal(err)
	}
	var owner = new(Owner)
	err = stmt.QueryRow(ownerEmail).Scan(&owner.Name, &owner.Email)
	if err != nil {
		log.Println("No email found")
	}
	return *owner, err
}

func SaveOwner (dbh sql.DB, o Owner) (int64, error) {
	stmt, err := dbh.Prepare("insert into owner (owner_name, owner_email) values (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(&o.Name, &o.Email)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return lastId, err
}