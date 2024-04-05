package db

import "log"

func CreateTables() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		team_id VARCHAR(255) NOT NULL,
		karma INT NOT NULL
	)`)

	if err != nil {
		log.Fatal("Error creating users table: ", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS groups (
		id SERIAL PRIMARY KEY,
		group_id VARCHAR(255) NOT NULL,
		team_id VARCHAR(255) NOT NULL,
		karma INT NOT NULL
	)`)

	if err != nil {
		log.Fatal("Error creating groups table: ", err)
	}
}
