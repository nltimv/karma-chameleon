package slack

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

var (
	webApi            *slack.Client
	socketMode        *socketmode.Client
	socketModeHandler *socketmode.SocketmodeHandler
	selfUserId        string
)

func Init(appToken string, botToken string, debugMode bool) {
	webApi = slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
		slack.OptionDebug(debugMode),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)

	socketMode = socketmode.New(
		webApi,
		socketmode.OptionDebug(debugMode),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)

	authTest, authTestErr := webApi.AuthTest()
	if authTestErr != nil {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
		os.Exit(1)
	}
	selfUserId = authTest.UserID
	fmt.Println("Authenticated successfully!")

	socketModeHandler = socketmode.NewSocketmodeHandler(socketMode)
}

func Run() {
	if socketMode == nil || socketModeHandler == nil {
		fmt.Fprintf(os.Stderr, "SocketMode is not initialized!\n")
		os.Exit(1)
	}

	socketModeHandler.RunEventLoop()
}

func AddEventHandler(eventHandler func(*socketmode.SocketmodeHandler)) {
	eventHandler(socketModeHandler)
}

func IsSelf(userId string) bool {
	return userId == selfUserId
}
