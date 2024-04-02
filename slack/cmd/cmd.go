package cmd

import (
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"nltimv.com/karma-chameleon/slack/internal/db"
	"nltimv.com/karma-chameleon/slack/internal/events"
	"nltimv.com/karma-chameleon/slack/internal/slack"
)

func Start() {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Slack bot token not provided")
	}

	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		log.Fatal("Slack app token not provided")
	}
	debugMode, _ := strconv.ParseBool(os.Getenv("SLACK_DEBUG_MODE"))
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	db.Open(dbHost, dbPort, dbUser, dbPassword, dbName)
	defer db.Close()

	db.CreateTables()

	slack.Init(appToken, botToken, debugMode)

	slack.AddEventHandler(events.HandleEvents)

	slack.Run()
}
