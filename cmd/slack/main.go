package main

import (
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"nltimv.com/karma-chameleon/internal/db"
	"nltimv.com/karma-chameleon/internal/events"
	"nltimv.com/karma-chameleon/internal/log"
	"nltimv.com/karma-chameleon/internal/slack"
)

func main() {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		log.Error.Fatal("Slack bot token not provided")
	}

	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		log.Error.Fatal("Slack app token not provided")
	}
	debugMode, _ := strconv.ParseBool(os.Getenv("SLACK_DEBUG_MODE"))
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSslMode := os.Getenv("DB_SSLMODE")

	db.Open(dbHost, dbPort, dbUser, dbPassword, dbName, dbSslMode)
	defer db.Close()

	db.Migrate()

	slack.Init(appToken, botToken, debugMode)

	slack.AddEventHandler(events.HandleEvents)

	slack.Run()
}
