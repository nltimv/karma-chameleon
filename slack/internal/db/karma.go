package db

import (
	"database/sql"
	"log"
)

func UpdateUserKarma(userID string, teamID string, increment int) (int, error) {

	var karma int
	var err error
	row := db.QueryRow("SELECT karma FROM users WHERE user_id = $1 AND team_id = $2", userID, teamID)
	if err = row.Scan(&karma); err != nil {
		if err == sql.ErrNoRows {
			karma = increment
			_, err := db.Exec("INSERT INTO users (user_id, team_id, karma) VALUES ($1, $2, $3)", userID, teamID, increment)
			if err != nil {
				log.Println("Error inserting user karma: ", err)
				return 0, err
			}
		} else {
			log.Println("Error scanning user karma: ", err)
			return 0, err
		}
	} else {
		karma += increment
		_, err := db.Exec("UPDATE users SET karma = $1 WHERE user_id = $2 AND team_id = $3", karma, userID, teamID)
		if err != nil {
			log.Println("Error updating user karma: ", err)
			return 0, err
		}
	}

	return karma, nil
}

func UpdateGroupKarma(groupID, teamID string, increment int) int {
	var karma int
	row := db.QueryRow("SELECT karma FROM groups WHERE group_id = $1 AND team_id = $2", groupID, teamID)
	if err := row.Scan(&karma); err != nil {
		if err == sql.ErrNoRows {
			karma = increment
			_, err := db.Exec("INSERT INTO groups (group_id, team_id, karma) VALUES ($1, $2, $3)", groupID, teamID, increment)
			if err != nil {
				log.Println("Error inserting group karma: ", err)
				return 0
			}
		} else {
			log.Println("Error scanning group karma: ", err)
			return 0
		}
	} else {
		karma += increment
		_, err := db.Exec("UPDATE groups SET karma = $1 WHERE group_id = $2 AND team_id = $3", karma, groupID, teamID)
		if err != nil {
			log.Println("Error updating group karma: ", err)
			return 0
		}
	}

	return karma
}
