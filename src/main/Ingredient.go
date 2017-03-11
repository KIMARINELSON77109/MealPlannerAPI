package main
// Fred T. Dunaway
// fred.t.dunaway@gmail.com
// March 3, 2017

import (
	"log"
	"database/sql"
	"strconv"
//	"time"

)

type Ingredient struct {
	Id			int		`json:"id"`
	Item		string  `json:"item"`
	Quantity	string	`json:"qunatity"`
	Brand		string	`json:"brand"`
	Uom			string	`json:"uom"`
	Usda_ndbno	string	`json:"usda_ndbno"`
}

func GetIngredient(dbh sql.DB, ingredient_id int) (Ingredient, error) {
	// finds the ingredient by the ingredient id
	stmt, err := dbh.Prepare("select item, quantity, brand, uom, usda_ndbno from ingredient, uom where ingredient.uom_id=uom.uom_id and ingredient.ingredient_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	var brand sql.NullString
	var usdaNdbno sql.NullString
	var quan float64
	var ing = new(Ingredient)
	defer stmt.Close()
	err = stmt.QueryRow(ingredient_id).Scan(&ing.Item, &quan, &brand, &ing.Uom, &usdaNdbno)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error finding ingredient with error: %v\n", err)
//			log.Fatal(err)
		}
	}
	if err != sql.ErrNoRows {
		if brand.Valid {
			ing.Brand = brand.String
		} else {
			ing.Brand = "Generic"
		}
		if usdaNdbno.Valid {
			ing.Usda_ndbno = usdaNdbno.String
		}
		ing.Quantity = strconv.FormatFloat(quan, 'f', 2, 32)
		if err != nil {
			log.Println("Error converting quantity")
			log.Fatal(err)
		}
		ing.Id = ingredient_id
	}
	stmt.Close()
	return *ing, nil
}

func SaveIngredient(dbh sql.DB, ingr Ingredient) (int64, error) {
	//BUG:  if no primary key is sent & the record is actually a duplicate, it saves it anyway
	query := "insert into ingredient (ingredient_id, item, quantity, brand, uom_id, usda_ndbno) " +
		"values (?, ?, ?, ?, (select uom_id from uom where uom=?), ?) " +
		"on duplicate key update ingredient_id=?, item=?, quantity=?, brand=?, uom_id=(select uom_id from uom where uom=?), usda_ndbno=?"
	stmt, err := dbh.Prepare(query)
	if err != nil {
		log.Println("Problem with save ingredient query")
		log.Println(query)
		log.Fatal(err)
	}
	qnty, qntConvertErr := strconv.ParseFloat(ingr.Quantity, 64)
	if qntConvertErr != nil {
		log.Println("Problem converting ingredient quantity")
		log.Fatal(err)
	}
//	log.Println("Dumping ingredient")
//	log.Printf("%+v\n", ingr)
	brand := ToNullString(ingr.Brand)
	usdaNdbno := ToNullString(ingr.Usda_ndbno)
	tx, err := dbh.Begin()
	if err != nil {
		log.Println("Can't start transaction for ingredient?")
		log.Fatal(err)
	}
	defer tx.Rollback()
	res, err := stmt.Exec(&ingr.Id, &ingr.Item, &qnty, &brand, &ingr.Uom, &usdaNdbno, &ingr.Id, &ingr.Item, &qnty, &brand, &ingr.Uom, &usdaNdbno)
	if err != nil {
		log.Printf("Problem executing ingredient save query with error: %v\n", err)
//		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Printf("Problem getting the last inserted id with error: %v\n", err)
//		log.Fatal(err)
	}
	stmt.Close()
	return lastId, err
}

func GetIngredientsForRecipe(dbh sql.DB, recipeId int) ([]Ingredient, error) {
	stmt, err := dbh.Prepare("SELECT ingredient.ingredient_id, `ingredient`.`item`, `ingredient`.`quantity`, `ingredient`.`brand`, `uom`.`uom`, usda_ndbno FROM `mooseware_meal_planner`.`ingredient_list` AS `ingredient_list`, `mooseware_meal_planner`.`ingredient` AS `ingredient`, `mooseware_meal_planner`.`recipe` AS `recipe`, `mooseware_meal_planner`.`uom` AS `uom` WHERE `ingredient_list`.`ingredient_id` = `ingredient`.`ingredient_id` AND `ingredient_list`.`recipe_id` = `recipe`.`recipe_id` AND `ingredient`.`uom_id` = `uom`.`uom_id` AND `recipe`.`recipe_id` = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(recipeId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var ingredients []Ingredient
	for rows.Next() {
		ingr := new(Ingredient)
		var brand sql.NullString
		var usda_ndbno sql.NullString
		err := rows.Scan(&ingr.Id, &ingr.Item, &ingr.Quantity, &brand, &ingr.Uom, &usda_ndbno)
		if err != nil {
			log.Println("get row failed with error: " + err.Error())
			return ingredients, err
		}
		if brand.Valid {
			ingr.Brand = brand.String
		} else {
			ingr.Brand = "Generic"
		}
		if usda_ndbno.Valid {
			ingr.Usda_ndbno = usda_ndbno.String
		}
		ingredients = append(ingredients, *ingr)
	}
	stmt.Close()
	return ingredients, err
}

