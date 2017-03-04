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
	Item		string  `json:"item"`
	Quantity	string	`json:"qunatity"`
	Brand		string	`json:"brand"`
	Uom			string	`json:"uom"`
}

func GetIngredient(dbh sql.DB, ingredient_id int) (Ingredient, error) {
	// finds the ingredient by the ingredient id
	stmt, err := dbh.Prepare("select item, quantity, brand, uom from ingredient, uom where ingredient.uom_id=uom.uom_id and ingredient.ingredient_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	var brand sql.NullString
	var quan float64
	var ing = new(Ingredient)
	err = stmt.QueryRow(ingredient_id).Scan(&ing.Item, &quan, &brand, &ing.Uom)
	if brand.Valid {
		ing.Brand = brand.String
	} else {
		ing.Brand = "Generic"
	}
	ing.Quantity = strconv.FormatFloat(quan, 'f', 2, 32)
	if err != nil {
		log.Fatal(err)
	}
	return *ing, nil
}

func SaveIngredient(dbh sql.DB, ingr Ingredient) (int64, error) {
	stmt, err := dbh.Prepare("insert into ingredient (item, quantity, brand, uom_id) values (?, ?, ?, (select uom_id from uom where uom=?))")
	if err != nil {
		log.Fatal(err)
	}
	qnty, _ := strconv.ParseFloat(ingr.Quantity, 64)
	res, err := stmt.Exec(&ingr.Item, &qnty, &ingr.Brand, &ingr.Uom)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return lastId, err
}

func GetIngredientsForRecipe(dbh sql.DB, recipeId int) ([]Ingredient, error) {
	stmt, err := dbh.Prepare("SELECT `ingredient`.`item`, `ingredient`.`quantity`, `ingredient`.`brand`, `uom`.`uom` FROM `mooseware_meal_planner`.`ingredient_list` AS `ingredient_list`, `mooseware_meal_planner`.`ingredient` AS `ingredient`, `mooseware_meal_planner`.`recipe` AS `recipe`, `mooseware_meal_planner`.`uom` AS `uom` WHERE `ingredient_list`.`ingredient_id` = `ingredient`.`ingredient_id` AND `ingredient_list`.`recipe_id` = `recipe`.`recipe_id` AND `ingredient`.`uom_id` = `uom`.`uom_id` AND `recipe`.`recipe_id` = ?")
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
		//TODO:  brand may be null.. deal with it.
		var brand sql.NullString
		err := rows.Scan(&ingr.Item, &ingr.Quantity, &brand, &ingr.Uom)
		if err != nil {
			log.Println("get row failed with error: " + err.Error())
			return ingredients, err
		}
		if brand.Valid {
			ingr.Brand = brand.String
		} else {
			ingr.Brand = "Generic"
		}
		ingredients = append(ingredients, *ingr)
	}
	return ingredients, err
}

