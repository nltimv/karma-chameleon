package slack

import (
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
		log.Fatalf("SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
	}
	selfUserId = authTest.UserID
	log.Println("Authenticated successfully!")

	socketModeHandler = socketmode.NewSocketmodeHandler(socketMode)
}

func Run() {
	if socketMode == nil || socketModeHandler == nil {
		log.Fatalln("SocketMode is not initialized!")
	}

	socketModeHandler.RunEventLoop()
}

func AddEventHandler(eventHandler func(*socketmode.SocketmodeHandler)) {
	eventHandler(socketModeHandler)
}

func IsSelf(userId string) bool {
	return userId == selfUserId
}
