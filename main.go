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
// main.go

package main

import "os"

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	a.Run(":8010")
}
