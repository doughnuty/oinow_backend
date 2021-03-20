package main

import (
	"database/sql"
	"sort"
)

type user struct {
	ID      int     `json:"id"`
	Name    string  `json:"first_name"`
	Surname string  `json:"last_name"`
	AituID  string  `json:"aituID"`
	Score   float64 `json:"score"`
	Style   int     `json:"style"`
	Phone   string  `json:"phone"`
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

func (game *games) getGameID(db *sql.DB) error {
	return db.QueryRow("SELECT id FROM games WHERE name=$1", game.Name).Scan(&game.ID)
}

func (u *user) createUserProfile(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO users (aituID, name, surname, phone, style) SELECT $1, $2, $3, $4, 0 WHERE NOT EXISTS (SELECT id FROM users WHERE aituID=$1)",
		u.AituID, u.Name, u.Surname, u.Phone)
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
	rows, err := db.Query("SELECT score FROM game_scores WHERE userID=$1", u.ID)
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
	rows, err := db.Query("SELECT id, name, surname, aituid FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []user{}

	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.Name, &u.Surname, &u.AituID); err != nil {
			return nil, err
		}

		if err := u.getScores(db); err != nil {
			return nil, err
		}

		users = append(users, u)
		sort.Slice(users, func(i, j int) bool {
			return users[i].Score < users[j].Score
		})
	}

	return users, nil
}

func (u *user) getScoreFromContacts(db *sql.DB) error {
	rows, err := db.Query("SELECT score FROM game_scores LEFT JOIN users ON userid=users.id WHERE phone=$1", u.Phone)
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
