package main

// Fred T. Dunaway
// fred.t.dunaway@gmail.com
// March 3, 2017

import (
	"log"
	"database/sql"
//	"strconv"
)

//TODO: add statemet close and only log fatal when really needed.
type Recipe struct {
	Id			int	`json:"id"`
	Name 		string `json:"name"`
	Directions	string `json:"directions"`
	Servings	int	`json:"servings"`
	OwnerName	string `json:"ownerName"`
	OwnerEmail	string	`json:"ownerEmail"`
	Ingredients[] Ingredient `json:"ingredients"`
}

func GetRecipeNameContains(dbh sql.DB, partialRecipeName string) ([]Recipe, error) {
	q := "select recipe_id, recipe_name, directions, servings, owner_name, owner_email from recipe, owner " +
		"where recipe_name like ?"
	stmt, err := dbh.Prepare(q)
	if err != nil {
		log.Fatalf("Unable to prepare query for partial recipe name match with error %v\n", err)
	}
	rows, err := stmt.Query("%" + partialRecipeName + "%")
	if err != nil {
		log.Printf("Query to find partial recipe name failed with error: %v\n", err)
		return nil, err
	}
	var recipes []Recipe
	for rows.Next() {
		rec := new(Recipe)
		var ownerName sql.NullString
		err := rows.Scan(&rec.Id, &rec.Name, &rec.Directions, &rec.Servings, &ownerName, &rec.OwnerEmail)
		if err != nil {
			log.Printf("Error retriving recipe containing %s with error: %v\n", partialRecipeName, err)
			return recipes, err
		}
		if ownerName.Valid {
			rec.OwnerName = ownerName.String
		}			
		recipes = append(recipes, *rec)
	}
	rows.Close()
	stmt.Close()
	return recipes, err
}

func GetRecipeNameSoundsLike(dbh sql.DB, nameSoundsLike string) ([]Recipe, error) {
		q := "select recipe_id, recipe_name, directions, servings, owner_name, owner_email from recipe, owner " +
			"where recipe.owner_id=owner.owner_id and SOUNDEX(?) = SOUNDEX(recipe.recipe_name) limit 100"
		stmt, err := dbh.Prepare(q)
		if err != nil {
			log.Fatalf("Uable to prepare query to find recipes that sound like with error: %v\n", err)
		}
		var recipes []Recipe
		rows, err := stmt.Query(nameSoundsLike)
		if err != nil {
			log.Printf("Soundex query for recipe name failed with error: %v\n", err)
			return nil, err
		}
		for rows.Next() {
			rec := new(Recipe)
			var ownerName sql.NullString
			err := rows.Scan(&rec.Id, &rec.Name, &rec.Directions, &rec.Servings, &ownerName, &rec.OwnerEmail)
			if err != nil {
				log.Printf("Error retriving recipe with name that sounds like %s with error: %v\n", nameSoundsLike, err)
				return recipes, err
			}
			if ownerName.Valid {
				rec.OwnerName = ownerName.String
			}			
			recipes = append(recipes, *rec)
		}
		rows.Close()
		stmt.Close()
		return recipes, err
	}

func GetRecipe(dbh sql.DB, recipeId int) (Recipe, error) {
	stmt, err := dbh.Prepare("select recipe_id, recipe_name, directions, servings, owner_name, owner_email from recipe, owner where recipe.owner_id=owner.owner_id and recipe_id=?")
		if err != nil {
		log.Fatal(err)
	}
	var recipe = new(Recipe)
	var ownerName sql.NullString
	err = stmt.QueryRow(recipeId).Scan(&recipe.Id, &recipe.Name, &recipe.Directions, &recipe.Servings, &ownerName, &recipe.OwnerEmail)
	if err != nil {
		log.Printf("Query for recipe retuned error: %v\n", err)
	}
	if ownerName.Valid {
		recipe.OwnerName = ownerName.String
	}
	stmt.Close()
	// now get the ingredients for this recipe
	recipe.Ingredients, err = GetIngredientsForRecipe(dbh, recipeId)
	return *recipe, err
}

func SaveRecipe(dbh sql.DB, recipe Recipe) (int64, error) {
	query := "INSERT INTO recipe (recipe_id, recipe_name, directions, servings, owner_id) " +
			" VALUES(?, ?, ?, ?, (select owner_id from owner where owner_email=?)) " +
			" on duplicate key update recipe_id=?, recipe_name=?, directions=?, servings=?, owner_id=(select owner_id from owner where owner_email=?)"
	stmt, err := dbh.Prepare(query)
	if err != nil {
		log.Println("Unable to prepare query.")
		log.Println(query)
		log.Fatal(err)
	}
	res, err := stmt.Exec(recipe.Id, recipe.Name, recipe.Directions, recipe.Servings, recipe.OwnerEmail, recipe.Id, recipe.Name, recipe.Directions, recipe.Servings, recipe.OwnerEmail)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	
	var ingredients = recipe.Ingredients
	for ingr := range ingredients {
//		log.Printf("Sending ingredient: %+v\n", ingredients[ingr])
		_, err := SaveIngredient(dbh, ingredients[ingr])
		if err != nil {
			log.Printf("Error saving ingredient %d\n", ingr)
		} 
	}
	stmt.Close()
	//update the ingredient_list table using a brute force method
	if lastId != 0 {
		txt, err := dbh.Begin()
		defer txt.Rollback()
		stmt, err = dbh.Prepare("delete from ingredient_list where recipe_id=?")
		if err != nil {
			log.Fatalf("Problem with query to delete recipe from ingredient_list table: %v\n", err)
		}
		_, err = stmt.Exec(lastId)
		if err != nil {
			log.Printf("Error deleting recipe from ingredient_list: %v\n", err)
		}
		stmt, err = dbh.Prepare("insert into ingredient_list (recipe_id, ingredient_id values (?, ?)")
		if err != nil {
			log.Fatalf("Unable to begin ingredient_list transaction: %v\n", err)
		}
		for ingr := range ingredients {
			log.Printf("Add ingredient %d to recipe %d\n", lastId, ingredients[ingr].Id)
			_, err = stmt.Exec(lastId, ingredients[ingr].Id)
			if err != nil {
				log.Printf("Error saving ingredient_list: %v\n", err)
			}
		}
		txt.Commit()
		stmt.Close()
	}
	return lastId, err	
}
