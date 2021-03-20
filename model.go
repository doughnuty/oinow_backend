package main

import (
	"database/sql"
)

type user struct {
	ID    int    `json:"id"`
	name  string `json:"name"`
	Phone string `json:"phone"`
}

type game_scores struct {
	ID     int `json:"id"`
	gameID int `json:"id"`
	userID int `json:"userid"`
	score  int `json:"score"`
}

type game struct {
	ID   int    `json:"id"`
	name string `json:"name"`
}

func (u *user) getUserID(db *sql.DB) error {
	return db.QueryRow("SELECT id FROM users WHERE phone=$1", u.Phone).Scan(&u.ID)
}

func (g *game_scores) updateUserScore(db *sql.DB) error {
	_, err := db.Exec("UPDATE game_scores SET score=$1 WHERE userID=$2 AND gameID=$3", g.score, g.userID, g.gameID)

	return err
}

func (u *user) createUserProfile(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO users(phone) VALUES($1) RETURNING id",
		u.Phone).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func getLeaderboard(db *sql.DB) ([]user, error) {
	rows, err := db.Query("SELECT name FROM users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []user{}

	for rows.Next() {
		var u user
		if err := rows.Scan(&u.name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
