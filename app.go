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

/*
	GET			SEND TO USER		PROCESS
	leaderboard	(name and score)	select from db
	score		(profile)			select from db
	room_id 	(id)				generate | add to db | send to user
*/
/*
	POST 			SEND TO BACK					RECEIVE				PROCESS
	game_results	(user_phone, game_name, score)	success
	login_user		(phone)							success				if user not in db register as new
	friends			(contacts)						names and scores
*/
func (a *App) handleRequests() {
	a.Router.HandleFunc("/rest/oinaw/profile/{aituID}", a.getUserScore).Methods("GET")
	a.Router.HandleFunc("/rest/oinaw/profile/", a.sendUserID).Methods("POST")
	a.Router.HandleFunc("/rest/oinaw/games/", a.createGame).Methods("POST")
}

func (a *App) getUserScore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	aituID := vars["aituID"]

	u := user{AituID: aituID}
	if err := u.getUserID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	err := a.DB.QueryRow("SELECT score FROM game_scores WHERE userID=$1 AND gameID=1", u.ID).Scan(&u.Score)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u.Score)
}

func (a *App) sendUserID(w http.ResponseWriter, r *http.Request) {
	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	fmt.Println("creating profile with user", u.ID, "|", u.Name, "|", u.AituID, "|", u.Score)
	if err := u.createUserProfile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// what is best to send?
	respondWithJSON(w, http.StatusCreated, u.Score)
}

func (a *App) createGame(w http.ResponseWriter, r *http.Request) {
	var g games
	decoder := json.NewDecoder(r.Body)
	fmt.Println("creating game1", g.Name)
	if err := decoder.Decode(&g); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	fmt.Println("creating game", g.Name)
	if err := g.createGameConfig(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// what is best to send?
	respondWithJSON(w, http.StatusCreated, g.ID)
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
