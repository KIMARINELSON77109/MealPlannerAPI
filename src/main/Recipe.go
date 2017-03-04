package main

// Fred T. Dunaway
// fred.t.dunaway@gmail.com
// March 3, 2017

import (
	"log"
	"database/sql"
//	"strconv"
)

type Recipe struct {
	Id			int	`json:"id"`
	Name 		string `json:"name"`
	Directions	string `json:"directions"`
	Servings	int	`json:"servings"`
	OwnerName	string `json:"ownerName"`
	OwnerEmail	string	`json:"ownerEmail"`
	Ingredients[] Ingredient `json:"ingredients"`
}

func GetRecipe(dbh sql.DB, recipeId int) (Recipe, error) {
	stmt, err := dbh.Prepare("select recipe_id, recipe_name, directions, servings, owner_name, owner_email from recipe, owner where recipe.owner_id=owner.owner_id and recipe_id=?")
		if err != nil {
		log.Fatal(err)
	}
	var recipe = new(Recipe)
	err = stmt.QueryRow(recipeId).Scan(&recipe.Id, &recipe.Name, &recipe.Directions, &recipe.Servings, &recipe.OwnerName, &recipe.OwnerEmail)
	if err != nil {
		log.Fatal(err)
	}
	// now get the ingredients for this recipe
	recipe.Ingredients, err = GetIngredientsForRecipe(dbh, recipeId)
	return *recipe, err
}