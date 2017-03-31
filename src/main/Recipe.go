package main

// Fred T. Dunaway
// fred.t.dunaway@gmail.com
// March 3, 2017

import (
	"log"
	"database/sql"
	"encoding/base64"
//	"strconv"
)

//TODO: add statemet close and only log fatal when really needed.
type Recipe struct {
	Id				int	`json:"id"`
	Name 			string `json:"name"`
	Directions		string `json:"directions"`
	Servings		int	`json:"servings"`
	RecipeOwner		Owner `json:"owner"`
	Picture			string	`json:"picture"`
	Ingredients[] Ingredient `json:"ingredients"`
}

func GetRecipeNameContains(dbh sql.DB, partialRecipeName string) ([]Recipe, error) {
	q := "select recipe_id, recipe_name, directions, servings, owner_id, picture from recipe, owner " +
		"where recipe_name like ?"
	stmt, err := dbh.Prepare(q)
	defer stmt.Close()
	if err != nil {
		log.Fatalf("Unable to prepare query for partial recipe name match with error %v\n", err)
	}
	rows, err := stmt.Query("%" + partialRecipeName + "%")
	defer rows.Close()
	if err != nil {
		log.Printf("Query to find partial recipe name failed with error: %v\n", err)
		return nil, err
	}
	var recipes []Recipe
	for rows.Next() {
		rec := new(Recipe)
		var recPic sql.NullString
		var ownerId int64
		err := rows.Scan(&rec.Id, &rec.Name, &rec.Directions, &rec.Servings, &ownerId, &recPic)
		if err != nil {
			log.Printf("Error retriving recipe containing %s with error: %v\n", partialRecipeName, err)
			return recipes, err
		}
		if recPic.Valid {
			picData :=[]byte(recPic.String)
			rec.Picture = base64.StdEncoding.EncodeToString(picData)
		}
		owner, err := GetOwnerById(dbh, ownerId)
		if err != nil {
			return recipes, err
		}
		rec.RecipeOwner = owner
		// now get the ingredients for this recipe
		rec.Ingredients, err = GetIngredientsForRecipe(dbh, rec.Id)						
		recipes = append(recipes, *rec)
	}
	rows.Close()
	stmt.Close()
	return recipes, err
}

func GetRecipeNameSoundsLike(dbh sql.DB, nameSoundsLike string) ([]Recipe, error) {
		q := "select recipe_id, recipe_name, directions, servings, owner_id, picture from recipe, owner " +
			"where recipe.owner_id=owner.owner_id and SOUNDEX(?) = SOUNDEX(recipe.recipe_name) limit 100"
		stmt, err := dbh.Prepare(q)
		defer stmt.Close()
		if err != nil {
			log.Fatalf("Unable to prepare query to find recipes that sound like with error: %v\n", err)
		}
		var recipes []Recipe
		rows, err := stmt.Query(nameSoundsLike)
		if err != nil {
			log.Printf("Soundex query for recipe name failed with error: %v\n", err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			rec := new(Recipe)
			var ownerId int64
			var pic sql.NullString
			err := rows.Scan(&rec.Id, &rec.Name, &rec.Directions, &rec.Servings, &ownerId, &pic)
			if err != nil {
				log.Printf("Error retriving recipe with name that sounds like %s with error: %v\n", nameSoundsLike, err)
				return recipes, err
			}
			if pic.Valid {
				picData :=[]byte(pic.String)
				rec.Picture = base64.StdEncoding.EncodeToString(picData)
			}
			rec.RecipeOwner, err = GetOwnerById(dbh, ownerId)
			if err != nil {
				return recipes, err
			} 
			// now get the ingredients for this recipe
			rec.Ingredients, err = GetIngredientsForRecipe(dbh, rec.Id)
			if err != nil {
				return recipes, err
			}
			recipes = append(recipes, *rec)
		}
		return recipes, err
	}

func GetRecipe(dbh sql.DB, recipeId int) (Recipe, error) {
	stmt, err := dbh.Prepare("select recipe_id, recipe_name, directions, servings, owner_id, picture from recipe where recipe_id=?")
		if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var recipe = new(Recipe)
	var ownerId int64
	var pic sql.NullString
	err = stmt.QueryRow(recipeId).Scan(&recipe.Id, &recipe.Name, &recipe.Directions, &recipe.Servings, &ownerId, &pic)
	if err != nil {
		log.Printf("Query for recipe retuned error: %v\n", err)
		return *recipe, err
	}
	if pic.Valid {
		picData :=[]byte(pic.String)
		recipe.Picture = base64.StdEncoding.EncodeToString(picData)
	}
	recipe.RecipeOwner, err = GetOwnerById(dbh, ownerId)
	if(err != nil) {
		return *recipe, err
	}
	// now get the ingredients for this recipe
	recipe.Ingredients, err = GetIngredientsForRecipe(dbh, recipeId)
	return *recipe, err
}

func SaveRecipe(dbh sql.DB, recipe Recipe) (int64, error) {
	query := "INSERT INTO recipe (recipe_id, recipe_name, directions, servings, owner_id, picture) " +
			" VALUES(?, ?, ?, ?, ?, ?) " +
			" on duplicate key update recipe_id=?, recipe_name=?, directions=?, servings=?, " + 
			" owner_id=?, picture=?"
	stmt, err := dbh.Prepare(query)
	defer stmt.Close()
	if err != nil {
		log.Println("Unable to prepare query.")
		log.Println(query)
		log.Fatal(err)
	}
	// decode picture string ...
	var picData []byte
	picData, err = base64.StdEncoding.DecodeString(recipe.Picture)
	if err != nil {
		log.Println("Error decoding picture")
	}
	var recId sql.NullInt64
	if recipe.Id != 0 {
		recId.Scan(recipe.Id)
	}
	res, err := stmt.Exec(recId, recipe.Name, recipe.Directions, recipe.Servings, recipe.RecipeOwner.Id, picData, recId, recipe.Name, recipe.Directions, recipe.Servings, recipe.RecipeOwner.Id, picData)
	if err != nil {
		log.Println("Error saving meal.")
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	var recipeId = int64(recipe.Id)
	if (recipeId != 0) && (recipeId != lastId) {
		lastId = recipeId
	}
	_, err = SaveOwner(dbh, recipe.RecipeOwner)
	if err != nil {
		return 0, err
	}
	var ingredients = recipe.Ingredients
	for ingr := range ingredients {
//		log.Printf("Sending ingredient: %+v\n", ingredients[ingr])
		_, err := SaveIngredient(dbh, ingredients[ingr])
		if err != nil {
			log.Printf("Error saving ingredient %d\n", ingr)
		} 
	}
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
		stmt.Close()
		stmt, err = dbh.Prepare("insert into ingredient_list (recipe_id, ingredient_id) values (?, ?)")
		if err != nil {
			log.Fatalf("Unable to begin ingredient_list transaction: %v\n", err)
		}
		for ingr := range ingredients {
			log.Printf("Add ingredient %d to recipe %d\n", ingredients[ingr].Id, lastId)
			_, err = stmt.Exec(lastId, ingredients[ingr].Id)
			if err != nil {
				log.Printf("Error saving ingredient list: %v\n", err)
				txt.Rollback()
				return 0, err
			}
		}
		txt.Commit()
		stmt.Close()
	}
	return lastId, err	
}
