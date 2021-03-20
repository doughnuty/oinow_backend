package main

import (
	"database/sql"
)

type user struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	AituID string  `json:"aituID"`
	Score  float64 `json:"score"`
}

type game_scores struct {
	ID     int `json:"id"`
	GameID int `json:"gameid"`
	UserID int `json:"userid"`
	Score  int `json:"score"`
}

type games struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (u *user) getUserID(db *sql.DB) error {
	return db.QueryRow("SELECT id FROM users WHERE aituID=$1", u.AituID).Scan(&u.ID)
}

func (g *game_scores) updateUserScore(db *sql.DB) error {
	_, err := db.Exec("UPDATE game_scores SET score=$1 WHERE userID=$2 AND gameID=$3", g.Score, g.UserID, g.GameID)

	return err
}

func (game *games) createGameConfig(db *sql.DB) error {
	_, _ = db.Exec("INSERT INTO games (name) VALUES($1)", game.Name)

	return nil
}

func (game *games) getGameScore(db *sql.DB) error {
	return db.QueryRow("SELECT id FROM games WHERE name=$1", game.Name).Scan(&game.ID)
}

func (u *user) createUserProfile(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO users (aituID, name) SELECT $1, $2 WHERE NOT EXISTS (SELECT id FROM users WHERE aituID=$1)", u.AituID, u.Name)
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT id FROM users WHERE aituID=$1", u.AituID).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *user) getScores(db *sql.DB) error {
	rows, err := db.Query("SELECT score FROM game_scores LEFT JOIN users ON game_scores.userID=users.id")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		score := 0.0
		if err := rows.Scan(&score); err != nil {
			return err
		}
		u.Score += score
	}

	return nil
}

func getLeaderboard(db *sql.DB) ([]user, error) {
	rows, err := db.Query("SELECT id, name, aituid FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []user{}

	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.Name, &u.AituID); err != nil {
			return nil, err
		}

		if err := u.getScores(db); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}
