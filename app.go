package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// hold application with db and router
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// connect to db
func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.handleRequests()
}

// start application
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) handleRequests() {
	a.Router.HandleFunc("/rest/oinaw/profile/{phone}", a.getUserScore).Methods("GET")
	// 	myRouter.HandleFunc("/articles", returnAllArticles)
	// 	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	// 	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	// 	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	// 	log.Fatal(http.ListenAndServe(":10000", myRouter))
	//
}

func (a *App) getUserScore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	phone := vars["phone"]

	u := user{Phone: phone}
	if err := u.getUserID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// score := 0
	// err := a.DB.QueryRow("SELECT score FROM game_scores WHERE userID=$1 AND gameID=1", u.ID).Scan(&score)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	respondWithJSON(w, http.StatusOK, u)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
