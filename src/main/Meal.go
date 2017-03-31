package main

import (
	"log"
	"database/sql"
)

type Meal struct {
	Id			int64	`json:"id"`
	Name 		string `json:"name"`
	Servings	int64  `json:"servings"`
	Account		Owner	`json:"owner"`
	Mealtype	string `json:"meal_type"`
	Profile		string `json:"meal_profile"`
	Recipies[]	Recipe	`json:"recipies"`
}

func GetMealsForProfile(dbh sql.DB, profile string) ([]Meal, error) {
	sql := "SELECT `meal`.`meal_id` " +
		"FROM `mooseware_meal_planner`.`meal_recipes` AS `meal_recipes`, `mooseware_meal_planner`.`meal` AS `meal`, " +
		"`mooseware_meal_planner`.`meal_profile` AS `meal_profile`, `mooseware_meal_planner`.`meal_type` AS `meal_type`, " +
		"`mooseware_meal_planner`.`owner` AS `owner`, `mooseware_meal_planner`.`recipe` AS `recipe` " +
		"WHERE `meal_recipes`.`meal_id` = `meal`.`meal_id` " +
		"AND `meal_recipes`.`meal_profile_id` = `meal_profile`.`meal_profile_id` " +
		"AND `meal`.`meal_type_id` = `meal_type`.`meal_type_id` AND `meal`.`account_id` = `owner`.`owner_id` " + 
		"AND `meal_recipes`.`recipe_id` = `recipe`.`recipe_id` AND `recipe`.`owner_id` = `owner`.`owner_id` " + 
		"AND `meal_profile`.`profile_name` = ?"
	stmt, err := dbh.Prepare(sql)
	if err != nil {
		log.Fatalf("Unable to prepare query for meals with profile:\n%v\n", err)
	}
	var meals []Meal
	rows, err := stmt.Query(profile)
	if err != nil {
		log.Printf("Unable to query for meals with profile:\n%v\n", err)
		return meals, err
	}
	for rows.Next() {
		var meal = new(Meal)
		var mealId int64
		err := rows.Scan(&mealId)
		if err != nil {
			log.Printf("Error retriving meal id for profile\n%v\n", err)
			return meals, err
		}
		*meal, err = GetMeal(dbh, mealId)
		if err != nil {
			return meals, err
		}
		meals = append(meals, *meal)
	}
	stmt.Close()
	return meals, err
}

func GetMeal(dbh sql.DB, mealId int64) (Meal, error) {
	sql := "SELECT `meal`.`meal_name`, `meal`.`servings`, `meal`.`account_id`, `meal_profile`.`profile_name`, `meal_type`.`meal_type` " +
		"FROM `mooseware_meal_planner`.`meal_recipes` AS `meal_recipes`, " +
		"`mooseware_meal_planner`.`meal` AS `meal`, `mooseware_meal_planner`.`meal_profile` AS `meal_profile`, " +
		"`mooseware_meal_planner`.`meal_type` AS `meal_type`, `mooseware_meal_planner`.`owner` AS `owner`, " +
		"`mooseware_meal_planner`.`recipe` AS `recipe` " + 
		"WHERE `meal_recipes`.`meal_id` = `meal`.`meal_id` " +
		"AND `meal_recipes`.`meal_profile_id` = `meal_profile`.`meal_profile_id` " +
		"AND `meal`.`meal_type_id` = `meal_type`.`meal_type_id` " +
		"AND `meal`.`account_id` = `owner`.`owner_id` " +
		"AND `meal_recipes`.`recipe_id` = `recipe`.`recipe_id` " +
		"AND `recipe`.`owner_id` = `owner`.`owner_id` AND `meal`.`meal_id` = ?"
	stmt, err := dbh.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		log.Println("Unable to create meal query")
		log.Fatal(err)
	}
	var meal = new(Meal)
	var accountId int64
	err = stmt.QueryRow(mealId).Scan(&meal.Name, &meal.Servings, &accountId, &meal.Profile, &meal.Mealtype)
	if err != nil {
		log.Println("Query for meal failed")
		log.Println(err)
		return *meal, err
	}
	stmt.Close()
	stmt, err = dbh.Prepare("select recipe_id from meal_recipes where meal_id=?")
	if err != nil {
		log.Println("Unable to create recipe for meal query")
		log.Fatal(err)
	}
	rows, err := stmt.Query(mealId)
	defer rows.Close()
	if err != nil {
		log.Println("Query for recipes for meal failed to execute")
		log.Println(err)
		return *meal, err
	}
	for rows.Next() {
		var recipeId int
		err := rows.Scan(&recipeId)
		if err != nil {
			log.Println("Unable to read recipe row")
			return *meal, err
		}
		recipe, err := GetRecipe(dbh, recipeId)
		if err != nil {
			return *meal, err
		}
		meal.Recipies = append(meal.Recipies, recipe)
	}
	meal.Id = mealId
	meal.Account, err = GetOwnerById(dbh, accountId)
	return *meal, err
}

func GetMealBySoundex(dbh sql.DB, soundsLike string) ([]Meal, error) {
	sql := "select meal_id from meal where SOUNDEX(?) = SOUNDEX(meal.meal_name) limit 100"
	stmt, err := dbh.Prepare(sql)
	if err != nil {
		log.Println("Error preparing meal soundex query")
		log.Fatal(err)
	}
	var meals []Meal
	rows, err := stmt.Query(soundsLike)
	defer rows.Close()
	defer stmt.Close()
	for rows.Next() {
		var mealId int64
		err := rows.Scan(&mealId)
		if err != nil {
			log.Println ("Error retriving meal in soundex query")
			log.Println(sql)
			return meals, err
		}
		m, err := GetMeal(dbh, mealId)
		if err != nil {
			return meals, err
		}
		meals = append(meals, m)
	}
	return meals, err
}

func SaveMeal(dbh sql.DB, meal Meal) (int64, error) {
	var query string
//	var newMealId int64
//	log.Printf("Meal dump: %+v\n", meal)
	if meal.Id != 0 {
		query = "insert into meal (meal_id, meal_name, servings, account_id, meal_type_id) " +
			"values (?, ?, ?, ?, " +
			"(select meal_type_id from meal_type where meal_type=?)) " +
			"on duplicate key update meal_name=?, servings=?, account_id=?, " +
			"meal_type_id=(select meal_type_id from meal_type where meal_type=?)"
	} else {
		query = "insert into meal (meal_name, servings, account_id, meal_type_id) " +
			"values (?, ?, ?, " +
			"(select meal_type_id from meal_type where meal_type=?)) "
	}
	txt, err := dbh.Begin()
	defer txt.Rollback()
	stmt, err := dbh.Prepare(query)
	defer stmt.Close()
	if err != nil {
		log.Printf("Unable to prepare meal save query: %s\n", query)
		log.Fatal(err)
	}
	var res sql.Result
	if meal.Id != 0 {
		res, err = stmt.Exec(meal.Id, meal.Name, meal.Servings, meal.Account.Id, meal.Mealtype, 
			meal.Name, meal.Servings, meal.Account.Id, meal.Mealtype)
	} else {
		res, err = stmt.Exec(meal.Name, meal.Servings, meal.Account.Id, meal.Mealtype)
	}
	if err != nil {
		log.Println("Unable to execute meal save query.")
		log.Println(err)
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Println("Unable to get the last meal id inserted")
		log.Println(err)
		return 0, err
	}
	if (meal.Id != 0) && (meal.Id != lastId) {
		lastId = meal.Id
	}
	stmt.Close()
	// update/insert to meal_recipes
	// use burte force method to deal w/ add/subtract of recipes from the list
	stmt, err = dbh.Prepare("delete from meal_recipes where meal_id=?")
	if err != nil {
		log.Fatalf("Unable to prepare query to delete from meal_recipes with error: %v\n", err)
	}
	_, err = stmt.Exec(lastId)
	if err != nil {
		log.Printf("Error deleting meal_recipes: %v\n", err)
	}
	stmt.Close()
	query = "insert into meal_recipes (meal_id, recipe_id, meal_profile_id) " +
		"values(?, ?, (select meal_type_id from meal_type where meal_type=?))"
	stmt, err = dbh.Prepare(query)
	if err != nil {
		log.Fatalf("Unable to prepare meal_recipes insert query: %v\n", err)
	}
	for r := range meal.Recipies {
		_, err = stmt.Exec(lastId, meal.Recipies[r].Id, meal.Mealtype)
		if err != nil {
			log.Printf("Unable to insert row into meal_recipes with error: %v\n", err)
			return 0, err
		}
		// deal with ingredients also
		for ing := range meal.Recipies[r].Ingredients {
			_, err := SaveIngredient(dbh, meal.Recipies[r].Ingredients[ing])
			if err != nil {
				log.Printf("Error saving ingredient for recipe: %v\n", err)
				return 0, err
			}
		}
	}
	stmt.Close()	
	txt.Commit()
	return lastId, err
}