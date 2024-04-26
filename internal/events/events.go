package events

import (
	"regexp"
	"strings"

	slackapi "github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"nltimv.com/karma-chameleon/internal/karma"
	"nltimv.com/karma-chameleon/internal/log"
	"nltimv.com/karma-chameleon/internal/slack"
	"nltimv.com/karma-chameleon/internal/ui"
)

const (
	interactionActionLeaderboardUsers = "actionLeaderboardUsers"
)

func HandleEvents(handler *socketmode.SocketmodeHandler) {
	handler.Handle(socketmode.EventTypeConnecting, handleConnecting)
	handler.Handle(socketmode.EventTypeConnectionError, handleConnectionError)
	handler.Handle(socketmode.EventTypeConnected, handleConnected)
	handler.Handle(socketmode.EventTypeHello, handleHello)
	handler.Handle(socketmode.EventTypeIncomingError, handleIncomingError)

	handler.HandleEvents(slackevents.Message, handleMessageEvent)
	handler.HandleEvents(slackevents.AppHomeOpened, handleAppHomeOpened)

	handler.HandleInteraction(slackapi.InteractionTypeBlockActions, handleInteraction)
}

func handleConnecting(evt *socketmode.Event, client *socketmode.Client) {
	log.Default.Println("Connecting to Slack with Socket Mode...")
}

func handleConnectionError(evt *socketmode.Event, client *socketmode.Client) {
	log.Default.Println("Connection failed. Retrying later...")
}

func handleConnected(evt *socketmode.Event, client *socketmode.Client) {
	log.Default.Println("Connected to Slack with Socket Mode.")
}

func handleHello(evt *socketmode.Event, client *socketmode.Client) {
	log.Default.Println("Hello from Slack!")
}

func handleIncomingError(evt *socketmode.Event, client *socketmode.Client) {
	log.Error.Println("Incoming error from Slack.")
}

func handleMessageEvent(evt *socketmode.Event, client *socketmode.Client) {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Default.Printf("Ignored %+v\n", evt)
		return
	}

	client.Ack(*evt.Request)

	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
	if !ok {
		log.Default.Printf("Ignored %+v\n", ev)
		return
	}

	if !slack.IsSelf(ev.User) {
		processMessage(ev, &eventsAPIEvent)
	}
}

func handleAppHomeOpened(evt *socketmode.Event, client *socketmode.Client) {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Default.Printf("Ignored %+v\n", evt)
		return
	}
	client.Ack(*evt.Request)

	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.AppHomeOpenedEvent)

	if !ok {
		log.Default.Printf("Ignored %+v\n", ev)
		return
	}

	if !slack.IsSelf(ev.User) {
		ui.ShowHomeTab(ev)
	}
}

func handleInteraction(evt *socketmode.Event, client *socketmode.Client) {
	interaction := evt.Data.(slackapi.InteractionCallback)
	client.Ack(*evt.Request)

	actionId := interaction.ActionCallback.BlockActions[0].ActionID

	if strings.HasPrefix(actionId, "actionLink") {
		return // no-op because the action is a link
	}

	switch actionId {
	case interactionActionLeaderboardUsers:
		karma.OpenLeaderboard(&interaction)
	default:
		log.Default.Printf("Unknown action ID: %s\n", actionId)
	}
}

func processMessage(ev *slackevents.MessageEvent, apiEvent *slackevents.EventsAPIEvent) {
	patterns := map[*regexp.Regexp]func(*slackevents.MessageEvent, *slackevents.EventsAPIEvent, *regexp.Regexp){
		regexp.MustCompile(`<@([a-zA-Z0-9_]+)>\s?(\+\+\+|---|\+\+|--)`):                              karma.ProcessUserKarma,
		regexp.MustCompile(`<!subteam\^([a-zA-Z0-9_]+)\|?[@a-zA-Z0-9_\-.]*>\s?(\+\+\+|---|\+\+|--)`): karma.ProcessGroupKarma,
		regexp.MustCompile(`<@([a-zA-Z0-9_]+)>\s?karma`):                                             karma.ProcessGetUserKarma,
		regexp.MustCompile(`<!subteam\^([a-zA-Z0-9_]+)\|?[@a-zA-Z0-9_\-.]*>\s?karma`):                karma.ProcessGetGroupKarma,
	}

	for pattern, processor := range patterns {
		if pattern.MatchString(ev.Text) {
			processor(ev, apiEvent, pattern)
			return
		}
	}
}
