package slack

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"

	"nltimv.com/karma-chameleon/internal/log"
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
		slack.OptionLog(log.Slack),
	)

	socketMode = socketmode.New(
		webApi,
		socketmode.OptionDebug(debugMode),
		socketmode.OptionLog(log.Slack),
	)

	authTest, authTestErr := webApi.AuthTest()
	if authTestErr != nil {
		log.Error.Fatalf("SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
	}
	selfUserId = authTest.UserID
	log.Default.Println("Authenticated successfully!")

	socketModeHandler = socketmode.NewSocketmodeHandler(socketMode)
}

func Run() {
	if socketMode == nil || socketModeHandler == nil {
		log.Default.Fatalln("SocketMode is not initialized!")
	}

	socketModeHandler.RunEventLoop()
}

func AddEventHandler(eventHandler func(*socketmode.SocketmodeHandler)) {
	eventHandler(socketModeHandler)
}

func IsSelf(userId string) bool {
	return userId == selfUserId
}

func GetAppUrl(teamId string, appId string, tab string) string {
	return fmt.Sprintf("slack://app?team=%s&id=%s&tab=%s", teamId, appId, tab)
}
