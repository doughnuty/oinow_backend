package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"

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
	leaderboard	(name and score)	select from db						DONE
	score		(profile)			select from db 						DONE
	room_id 	(id)				generate | add to db | send to user DONE
*/
/*
	POST 			SEND TO BACK					RECEIVE				PROCESS
	game_results	(user_id, game_name, score)		success												  DONE
	login_user		(user_id)						success				if user not in db register as new DONE
	friends			(contacts)						names and scores									  DONE
	create_game		(game name)						success				add game to db if absent          DONE
	buy_pretty		(style, price, aituID)			style and newscore	substract price from score and upd style
*/

func (a *App) handleRequests() {
	a.Router.HandleFunc("/rest/oinow/profile/{aituID}", a.getUserScore).Methods("GET")
	a.Router.HandleFunc("/rest/oinow/profile/", a.sendUserID).Methods("POST")
	a.Router.HandleFunc("/rest/oinow/games/", a.createGame).Methods("POST")
	a.Router.HandleFunc("/rest/oinow/leaderboard/", a.getLeaderboard).Methods("GET")
	a.Router.HandleFunc("/rest/oinow/friends/", a.getFriendsList).Methods("POST")
	a.Router.HandleFunc("/rest/oinow/new_game/", a.generateRoom).Methods("GET")
	a.Router.HandleFunc("/rest/oinow/profile/results/", a.getGameResults).Methods("POST")
	a.Router.HandleFunc("/rest/oinow/profile/shop", a.buyPretty).Methods("POST")
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

	err := u.getScores(a.DB)
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
	if err := u.createUserProfile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := u.getScores(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// what is best to send?
	respondWithJSON(w, http.StatusCreated, u.Score)
}

func (a *App) createGame(w http.ResponseWriter, r *http.Request) {
	var g games
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&g); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := g.createGameConfig(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// what is best to send?
	respondWithJSON(w, http.StatusCreated, `Success`)
}

func (a *App) getLeaderboard(w http.ResponseWriter, r *http.Request) {
	users, err := getLeaderboard(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) getFriendsList(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	users := []user{}

	if err := decoder.Decode(&users); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	for id := range users {
		if err := users[id].getScoreFromContacts(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Score > users[j].Score
	})
	respondWithJSON(w, http.StatusCreated, users)
}

func (a *App) generateRoom(w http.ResponseWriter, r *http.Request) {
	randStr := RandStringBytes(6)

	respondWithJSON(w, http.StatusCreated, randStr)
}

func (a *App) getGameResults(w http.ResponseWriter, r *http.Request) {

	type gameresults struct {
		User, Game_name string
		Score           int
	}

	var result gameresults
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&result); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	user := user{AituID: result.User}
	if err := user.getUserID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	game := games{Name: result.Game_name}
	if err := game.getGameID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	scores := game_scores{
		UserID: user.ID,
		GameID: game.ID,
		Score:  result.Score}

	if err := scores.updateUserScore(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, `Success`)
}

func (a *App) buyPretty(w http.ResponseWriter, r *http.Request) {
	type receipt struct {
		AituID string
		Price  float64
		Style  int
	}

	var rec receipt

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&rec); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	user := user{AituID: rec.AituID, Style: rec.Style}

	if err := user.UpdateStyle(a.DB, rec.Price); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := user.getUserbyID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
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
