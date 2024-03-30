package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"

	_ "github.com/lib/pq"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"

	"nltimv.com/karma-chameleon/slack/internal/events"
)

var (
	db *sql.DB
	sm *socketmode.Client
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

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	createTables()

	slackWebApi := slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
		slack.OptionDebug(false),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)

	sm = socketmode.New(
		slackWebApi,
		socketmode.OptionDebug(false),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)

	authTest, authTestErr := slackWebApi.AuthTest()
	if authTestErr != nil {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
		os.Exit(1)
	}
	selfUserId := authTest.UserID
	fmt.Printf("Authenticated successfully! User ID: %s\n", selfUserId)

	go events.HandleEvents(sm, slackWebApi, selfUserId)

	sm.Run()
}

func createTables() {
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

func processMessage(ev *slack.MessageEvent, rtm *slack.RTM) {
	reUserKarma := regexp.MustCompile(`<@([a-zA-Z0-9_]+)>\s?(\+\+\+|---|\+\+|--)`)
	reGroupKarma := regexp.MustCompile(`<!subteam\^([a-zA-Z0-9_]+)\|?[@a-zA-Z0-9_\-.]*>\s?(\+\+\+|---|\+\+|--)`)

	if reUserKarma.MatchString(ev.Text) {
		processUserKarma(ev, reUserKarma, rtm)
	} else if reGroupKarma.MatchString(ev.Text) {
		processGroupKarma(ev, reGroupKarma, rtm)
	}
}

func processUserKarma(ev *slack.MessageEvent, re *regexp.Regexp, rtm *slack.RTM) {
	matches := re.FindStringSubmatch(ev.Text)
	userID := matches[1]
	increment := matches[2]

	incrementValue := 0
	switch increment {
	case "+++":
		incrementValue = 2
	case "++":
		incrementValue = 1
	case "---":
		incrementValue = -2
	case "--":
		incrementValue = -1
	}

	if incrementValue == 0 {
		return
	}

	karma := updateUserKarma(userID, ev.Team, incrementValue)

	var response string
	switch incrementValue {
	case 2:
		response = fmt.Sprintf("<@%s>'s karma got a double boost! 🚀 New karma count: %d", userID, karma)
	case 1:
		response = fmt.Sprintf("<@%s>'s karma is on the rise! 🚀 New karma count: %d", userID, karma)
	case -1:
		response = fmt.Sprintf("<@%s>'s karma took a hit! 💔 New karma count: %d", userID, karma)
	case -2:
		response = fmt.Sprintf("<@%s>'s karma took a double hit! 💔 New karma count: %d", userID, karma)
	}

	rtm.SendMessage(rtm.NewOutgoingMessage(response, ev.Channel))
}

func updateUserKarma(userID, teamID string, increment int) int {
	if !isValidUser(userID) {
		return 0
	}

	var karma int
	row := db.QueryRow("SELECT karma FROM users WHERE user_id = $1 AND team_id = $2", userID, teamID)
	if err := row.Scan(&karma); err != nil {
		if err == sql.ErrNoRows {
			karma = increment
			_, err := db.Exec("INSERT INTO users (user_id, team_id, karma) VALUES ($1, $2, $3)", userID, teamID, increment)
			if err != nil {
				log.Println("Error inserting user karma: ", err)
				return 0
			}
		} else {
			log.Println("Error scanning user karma: ", err)
			return 0
		}
	} else {
		karma += increment
		_, err := db.Exec("UPDATE users SET karma = $1 WHERE user_id = $2 AND team_id = $3", karma, userID, teamID)
		if err != nil {
			log.Println("Error updating user karma: ", err)
			return 0
		}
	}

	return karma
}

func updateGroupKarma(groupID, teamID string, increment int) int {
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

func isValidUser(userID string) bool {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	userInfo, err := api.GetUserInfo(userID)
	if err != nil {
		log.Println("Error getting user info: ", err)
		return false
	}

	return userInfo != nil && !userInfo.Deleted
}

func processGroupKarma(ev *slack.MessageEvent, re *regexp.Regexp, rtm *slack.RTM) {
	matches := re.FindStringSubmatch(ev.Text)
	groupID := matches[1]
	increment := matches[2]

	incrementValue := 0
	switch increment {
	case "+++":
		incrementValue = 2
	case "++":
		incrementValue = 1
	case "---":
		incrementValue = -2
	case "--":
		incrementValue = -1
	}

	if incrementValue == 0 {
		return
	}

	usergroupMembers := getUsergroupMembers(groupID)

	if len(usergroupMembers) == 0 {
		return
	}

	for _, memberID := range usergroupMembers {
		if memberID != ev.User || incrementValue < 0 {
			_ = updateUserKarma(memberID, ev.Team, incrementValue)
		}
	}

	karma := updateGroupKarma(groupID, ev.Team, incrementValue)

	var response string
	switch incrementValue {
	case 2:
		response = fmt.Sprintf("The karma of <!subteam^%s> and its members got a double boost! 🚀 New group karma count: %d", groupID, karma)
	case 1:
		response = fmt.Sprintf("The karma of <!subteam^%s> and its members is on the rise! 🚀 New group karma count: %d", groupID, karma)
	case -1:
		response = fmt.Sprintf("The karma of <!subteam^%s> and its members took a hit! 💔 New group karma count: %d", groupID, karma)
	case -2:
		response = fmt.Sprintf("The karma of <!subteam^%s> and its members took a double hit! 💔 New group karma count: %d", groupID, karma)
	}

	rtm.SendMessage(rtm.NewOutgoingMessage(response, ev.Channel))
}

func getUsergroupMembers(groupID string) []string {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	usergroup, err := api.GetUserGroupMembers(groupID)
	if err != nil {
		log.Println("Error getting user group members: ", err)
		return []string{}
	}

	return usergroup
}
