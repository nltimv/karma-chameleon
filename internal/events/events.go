package events

import (
	"fmt"
	"regexp"

	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"nltimv.com/karma-chameleon/internal/karma"
	"nltimv.com/karma-chameleon/internal/slack"
)

func HandleEvents(handler *socketmode.SocketmodeHandler) {
	handler.Handle(socketmode.EventTypeConnecting, handleConnecting)
	handler.Handle(socketmode.EventTypeConnectionError, handleConnectionError)
	handler.Handle(socketmode.EventTypeConnected, handleConnected)
	handler.Handle(socketmode.EventTypeHello, handleHello)

	handler.HandleEvents(slackevents.Message, handleMessageEvent)
}

func handleConnecting(evt *socketmode.Event, client *socketmode.Client) {
	fmt.Println("Connecting to Slack with Socket Mode...")
}

func handleConnectionError(evt *socketmode.Event, client *socketmode.Client) {
	fmt.Println("Connection failed. Retrying later...")
}

func handleConnected(evt *socketmode.Event, client *socketmode.Client) {
	fmt.Println("Connected to Slack with Socket Mode.")
}

func handleHello(evt *socketmode.Event, client *socketmode.Client) {
	fmt.Println("Hello from Slack!")
}

func handleMessageEvent(evt *socketmode.Event, client *socketmode.Client) {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}

	client.Ack(*evt.Request)

	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", ev)
		return
	}

	if !slack.IsSelf(ev.User) {
		processMessage(ev, &eventsAPIEvent)
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
