package main
// Fred T. Dunaway
// fred.t.dunaway@gmail.com
// March 3, 2017

import (
   "encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
    "io"
	"log"
    "github.com/gorilla/mux"
	"strconv"
	"os"
)

const maxReadBytes = 1048576
const internalServerError = 500

func init() {
	file, configFileErr := os.Open("dbconfig.json")
	if configFileErr != nil {
		log.Panicln("Can not open dbconfig.json config file.")
	}
	decoder := json.NewDecoder(file)
	mydbp := DatabaseConnectionPrameters{}
	err := decoder.Decode(&mydbp)
	if err != nil {
		log.Println("error reading config: " + err.Error())
		log.Panicln("Unable to read database configuration file (dbconfig.json)")
	}
	log.Println("getting new database handler")
	dbh, err = NewDBH(mydbp)	//database hanlder, dbh, declared in DatabaseHelper as global
	if err != nil {
		log.Fatal("oops.... Database handler didn't initialize")
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

func IngredientGet (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    vars := mux.Vars(r)
    fmt.Printf("vars: %v\n", vars)
    ingredient_id := vars["ingredient_id"]
    if len(ingredient_id) != 0 {
    	ingr_id, _ := strconv.Atoi(ingredient_id)
    	ingredient, err := GetIngredient(*dbh, ingr_id)
     	if err != nil {
    		w.WriteHeader(http.StatusNotFound)
    		log.Println(err)
    	} else {
    		if err := json.NewEncoder(w).Encode(ingredient); err != nil {
				panic(err)
		    }
    	}
    } else {
    	w.WriteHeader(http.StatusBadRequest)
    }	
}

func IngredientCreate (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
	var ingr Ingredient
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxReadBytes))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &ingr); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }
	ingId, nfgErr := SaveIngredient(*dbh, ingr)
	if(nfgErr != nil) {
		log.Print(nfgErr)
		w.WriteHeader(422)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }		
	}
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(ingId); err != nil {
		panic(err)
    }
}

func OwnerCreate (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
	var owner Owner
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxReadBytes))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &owner); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }
	id, nfgErr := SaveOwner(*dbh, owner)
	if(nfgErr != nil) {
		log.Print(nfgErr)
		w.WriteHeader(422)
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }		
	}
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(id); err != nil {
		panic(err)
    }	
}

func OwnerGet (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    vars := mux.Vars(r)
    fmt.Printf("vars: %v\n", vars)
    ownerEmail := vars["ownerEmail"]
    if len(ownerEmail) != 0 {
    	owner, err := GetOwner(*dbh, ownerEmail)
     	if err != nil {
    		w.WriteHeader(http.StatusNotFound)
    		log.Println(err)
    	} else {
    		if err := json.NewEncoder(w).Encode(owner); err != nil {
				panic(err)
		    }
    	}
    } else {
    	w.WriteHeader(http.StatusBadRequest)
    }	
}

func RecipeGet (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    vars := mux.Vars(r)
    fmt.Printf("vars: %v\n", vars)
    recipieId := vars["recipe_id"]
    if len(recipieId) != 0 {
    	rId, err := strconv.Atoi(recipieId)
    	recipe, err := GetRecipe(*dbh, rId)
     	if err != nil {
    		w.WriteHeader(http.StatusNotFound)
    		log.Println(err)
    	} else {
    		if err := json.NewEncoder(w).Encode(recipe); err != nil {
				panic(err)
		    }
    	}
    } else {
    	w.WriteHeader(http.StatusBadRequest)
    }	
}
