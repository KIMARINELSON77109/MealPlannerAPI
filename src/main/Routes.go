package main
// Fred T. Dunaway
// fred.t.dunaway@gmail.com
// March 3, 2017

import (
	"net/http"
//	"database/sql"
)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
    Route{
        "Index",
        "GET",
        "/",
        Index,
    },
    Route{
    	"IngredientCreate",
    	"POST",
    	"/ingredient",
    	IngredientCreate,
    },
    Route{
    	"IngredientGet",
    	"GET",
    	"/ingredient/{ingredient_id:[0-9]+}",
    	IngredientGet,
    },
    Route{
    	"OwnerCreate",
    	"POST",
    	"/owner",
    	OwnerCreate,
    },
    Route{
    	"OwnerGet",
    	"GET",
    	"/owner/{ownerEmail}",
    	OwnerGet,
    },
    Route{
    	"RecipeGet",
    	"GET",
    	"/recipe/{recipe_id:[0-9]+}",
    	RecipeGet,
    },
    Route{
    	"RecipeCreate",
    	"POST",
    	"/recipe",
    	RecipeCreate,
    },
    Route{
    	"FindRecipeNameSoundsLike",
    	"GET",
    	"/recipe/soundslike/{recipeSoundsLikeName}",
    	FindRecipeNameSoundsLike,
    },
    Route{
    	"FindRecipeNameContains",
    	"GET",
    	"/recipe/contains/{recipeNameContains}",
    	FindRecipeNameContains,
    },
    Route{
    	"GetMealById",
    	"GET",
    	"/meal/{mealId:[0-9]+}",
    	GetMealById,
    },
    Route{
    	"GetMealsSoundsLike",
    	"GET",
    	"/meals/soundslike/{mealSoundsLike}",
    	GetMealsSoundsLike,
    },
    // no handler implemented yet
    Route{
    	"GetMealsForProfileHandler",
    	"GET",
    	"/meals/profile/{profileId}",
    	GetMealsForProfileHandler,
    },
    Route{
    	"SaveMealHandler",
    	"POST",
    	"/meal",
    	SaveMealHandler,
    },
}