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
    // everything below is for reference only.  It doesn't belong here.
/*    
    Route{
    	"PatientLocationManager",
    	"GET",
    	"/LocationManager/Patient/latitude/{lat}/longitude/{long}/floor/{floor}/LocationError/{locErr}",
    	PatientLocationManager,
    },
    Route{
    	"CreateCarePlan",
    	"POST",
    	"/CarePlan",
    	CreateCarePlan,   	
    },
    // The following is a dummy route used only for independent demo's.
    // It will be deprecated as soon as the video demo's using the IPS are running
    Route{
     	"NextPatient",
    	"GET",
    	"/demo/NextPatient/{patientId}",
    	NextPatient,   	
    },
    // this is only used for independent demo's w/o the IPS.
    // it will be deprecated
    Route {
    	"AddPatientToQue",
    	"POST",
    	"/demo/AddPatientToQue",
    	HandleAddPatientToQueue,
    },
    Route {
    	"CreateObservation",
    	"POST",
    	"/Observation",
    	HandleCreateObservation,
    },
    Route {
    	"GetObservations",
    	"GET",
    	"/Observations/{patient_id:[0-9]+}",
    	GetObservations,
    },
*/    
}