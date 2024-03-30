package events

import (
	"fmt"
	"log"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func HandleEvents(sm *socketmode.Client, slackWebApi *slack.Client, selfUserId string) {
	for ev := range sm.Events {
		switch ev.Type {
		case socketmode.EventTypeConnecting:
			fmt.Println("Connecting to Slack with Socket Mode...")
		case socketmode.EventTypeConnectionError:
			fmt.Println("Connection failed. Retrying later...")
		case socketmode.EventTypeConnected:
			fmt.Println("Connected to Slack with Socket Mode.")
		case socketmode.EventTypeEventsAPI:
			sm.Ack(*ev.Request)
			eventPayload, _ := ev.Data.(slackevents.EventsAPIEvent)

			switch eventPayload.Type {
			case string(slackevents.Message):
			case slackevents.CallbackEvent:
				switch event := eventPayload.InnerEvent.Data.(type) {
				case *slackevents.MessageEvent:
					if event.User != selfUserId &&
						strings.Contains(strings.ToLower(event.Text), "hello") {
						_, _, err := slackWebApi.PostMessage(
							event.Channel,
							slack.MsgOptionText(
								fmt.Sprintf(":wave: Hi there, <@%v>!", event.User),
								false,
							),
						)
						if err != nil {
							log.Printf("Failed to reply: %v", err)
						}
					}
				default:
					sm.Debugf("Skipped: %v", event)
				}
			default:
				sm.Debugf("unsupported Events API eventPayload received")
			}
		}
	}
}
