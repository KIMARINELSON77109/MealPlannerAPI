package main

import (
	"log"
	"database/sql"
//	"strconv"
)

type Owner struct {
	Id		int64	`json:"id"`
	Name 	string `json:"name"`
	Email	string `json:"email"`
}

func GetOwner (dbh sql.DB, ownerEmail string) (Owner, error){	
	stmt, err := dbh.Prepare("select owner_id, owner_name, owner_email from owner where owner_email=?")
	if err != nil {
		log.Fatal(err)
	}
	var owner = new(Owner)
	var oName sql.NullString
	err = stmt.QueryRow(ownerEmail).Scan(&owner.Id, &oName, &owner.Email)
	if err != nil {
		log.Println("No email found")
	}
	if oName.Valid {
		owner.Name = oName.String
	}
	stmt.Close()
	return *owner, err
}

func GetOwnerById (dbh sql.DB, ownerId int64) (Owner, error) {
	stmt, err := dbh.Prepare("select owner_name, owner_email from owner where owner_id=?")
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	var owner = new(Owner)
	var oName sql.NullString
	err = stmt.QueryRow(ownerId).Scan(&oName, &owner.Email)
	if err != nil {
		log.Printf("Unable to find owner with id: %d\n", ownerId)
		return *owner, err
	}
	owner.Id = ownerId
	if oName.Valid {
		owner.Name = oName.String
	}
	stmt.Close()
	return *owner, err
}

func SaveOwner (dbh sql.DB, o Owner) (int64, error) {
	q := "insert into owner (owner_id, owner_name, owner_email) values (?, ?, ?) " +
		"on duplicate key update owner_id=?, owner_name=?, owner_email=?"
	stmt, err := dbh.Prepare(q)
	if err != nil {
		log.Fatal(err)
	}
	//TODO:  deal with possile null's in id & name.
	name := ToNullString(o.Name)
	res, err := stmt.Exec(&o.Id, &name, &o.Email, &o.Id, &name, &o.Email)
	if err != nil {
		log.Printf("Error saving owner: %v\n",err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error getting the last id after saving owner: %v\n", err)
	}
	stmt.Close()
	return lastId, err
}