package slack

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

type ExtraFields struct {
	Label string
	Value string
}

type MessageContext struct {
	ChannelId string
	TeamId    string
	AppId     string
}

// func Say(message string, channel string) {
// 	_, _, err := webApi.PostMessage(channel, slack.MsgOptionText(message, false))
// 	if err != nil {
// 		fmt.Printf("failed posting message: %v", err)
// 	}
// }

func Say(message string, ctx *MessageContext, extraFields []*ExtraFields) {
	fields := make([]*slack.TextBlockObject, 0, len(extraFields))
	for _, extraField := range extraFields {

		fields = append(fields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*%s:*\n%s", extraField.Label, extraField.Value), false, false))
	}
	messageBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", message, false, false), fields, nil)
	divider := slack.NewDividerBlock()
	buttonBlock := getStandardButtons(ctx)

	_, _, err := webApi.PostMessage(ctx.ChannelId, slack.MsgOptionBlocks(messageBlock, divider, buttonBlock))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed posting message: %v\n", err)
	}
}

func getStandardButtons(ctx *MessageContext) *slack.ActionBlock {
	leaderboardButton := slack.NewButtonBlockElement("actionLeaderboardUsers", "valueLeaderboard", slack.NewTextBlockObject("plain_text", ":trophy: Leaderboard", true, false))
	helpButton := slack.NewButtonBlockElement("actionLinkHelp", "valueLinkHelp", slack.NewTextBlockObject("plain_text", ":question: Help / About", true, false)).WithURL(GetAppUrl(ctx.TeamId, ctx.AppId, "home"))
	return slack.NewActionBlock("msgShortcuts", leaderboardButton, helpButton)
}
