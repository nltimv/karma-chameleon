package slack

import (
	"github.com/slack-go/slack"

	"nltimv.com/karma-chameleon/internal/log"
)

func OpenModal(triggerId string, view slack.ModalViewRequest) {
	_, err := webApi.OpenView(triggerId, view)
	if err != nil {
		log.Error.Println("Error opening modal", err)
	}
}

func PublishHomeTab(userId string, view slack.HomeTabViewRequest, hash string) {
	_, err := webApi.PublishView(userId, view, hash)
	if err != nil {
		log.Error.Println("Error publishing home tab", err)
	}
}
