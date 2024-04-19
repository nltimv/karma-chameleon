package ui

import (
	"fmt"
	"os"

	slackapi "github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"nltimv.com/karma-chameleon/internal/slack"
)

func ShowHomeTab(ev *slackevents.AppHomeOpenedEvent) {
	blocks := make([]slackapi.Block, 0)
	var fields []*slackapi.TextBlockObject

	blocks = append(blocks, slackapi.NewHeaderBlock(
		slackapi.NewTextBlockObject("plain_text", "Welcome to the Karma Chameleon Home Page!", true, false),
	))

	blocks = append(blocks, slackapi.NewActionBlock(
		"blockActionButtons",
		slackapi.NewButtonBlockElement(
			"actionLeaderboardUsers",
			"valueLeaderboardUsers",
			slackapi.NewTextBlockObject("plain_text", ":medal: Show leaderboard (users)", true, false),
		),
	))

	blocks = append(blocks, slackapi.NewDividerBlock())

	blocks = append(blocks, slackapi.NewHeaderBlock(
		slackapi.NewTextBlockObject("plain_text", ":question: Help", true, false),
	))

	blocks = append(blocks, slackapi.NewSectionBlock(
		slackapi.NewTextBlockObject("plain_text", "To use the bot, you'll need to add it to a public or private channel in Slack. Then send one of the following messages to interact with the bot:", true, false), nil, nil))

	fields = []*slackapi.TextBlockObject{
		slackapi.NewTextBlockObject("mrkdwn", "`@user ++`", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "Awards one karma to `@user`. Does not work on yourself.", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "`@user +++`", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "Awards two karma to `@user`. Does not work on yourself.", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "`@user --`", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "Takes away one karma from `@user`.", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "`@user ---`", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "Takes away two karma from `@user`.", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "`@user karma`", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "Returns the current karma of `@user`.", false, false),
	}
	blocks = append(blocks, slackapi.NewSectionBlock(
		slackapi.NewTextBlockObject("mrkdwn", "*User interactions*", false, false),
		fields,
		nil,
	))

	fields = []*slackapi.TextBlockObject{
		slackapi.NewTextBlockObject("mrkdwn", "`@group ++`", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "Awards one karma to `@group` and all its members, excluding yourself. Does not work if you're the only member of the group.", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "`@group +++`", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "Awards two karma to `@group` and all its members, excluding yourself. Does not work if you're the only member of the group.", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "`@group --`", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "Takes away one karma from `@group` and all its members.", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "`@group ---`", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "Takes away two karma from `@group` and all its members.", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "`@group karma`", false, false),
		slackapi.NewTextBlockObject("mrkdwn", "Returns the current karma of `@group`.", false, false),
	}

	blocks = append(blocks, slackapi.NewSectionBlock(
		slackapi.NewTextBlockObject("mrkdwn", "*User group interactions*", false, false),
		fields,
		nil,
	))

	blocks = append(blocks, slackapi.NewDividerBlock())

	blocks = append(blocks, slackapi.NewHeaderBlock(
		slackapi.NewTextBlockObject("plain_text", ":information_source: About Karma Chameleon", true, false),
	))

	blocks = append(blocks, slackapi.NewSectionBlock(
		slackapi.NewTextBlockObject("plain_text", "© 2024 nltimv and contributors", true, false),
		nil, nil,
	))

	blocks = append(blocks, slackapi.NewSectionBlock(
		slackapi.NewTextBlockObject("plain_text", "This project is licensed under the MIT license.", true, false),
		nil, nil,
	))

	versionNumber := os.Getenv("APP_VERSION")

	blocks = append(blocks, slackapi.NewContextBlock(
		"blockVersion",
		slackapi.NewImageBlockElement(
			"https://raw.githubusercontent.com/nltimv/karma-chameleon/main/assets/img/logo.jpg",
			"cute cat",
		),
		slackapi.NewTextBlockObject("plain_text", fmt.Sprintf("Karma Chameleon - Version %s", versionNumber), true, false),
	))

	blocks = append(blocks, slackapi.NewActionBlock(
		"blockLinkButtons",
		slackapi.NewButtonBlockElement(
			"actionLinkGitHub",
			"valueLinkGitHub",
			slackapi.NewTextBlockObject("plain_text", ":globe_with_meridians: GitHub", true, false),
		).WithURL("https://github.com/nltimv/karma-chameleon"),
		slackapi.NewButtonBlockElement(
			"actionLinkLicense",
			"valueLinkLicense",
			slackapi.NewTextBlockObject("plain_text", ":page_with_curl: License", true, false),
		).WithURL("https://github.com/nltimv/karma-chameleon/blob/main/LICENSE"),
	))

	homeTab := slackapi.HomeTabViewRequest{
		Type:   slackapi.VTHomeTab,
		Blocks: slackapi.Blocks{BlockSet: blocks},
	}

	slack.PublishHomeTab(ev.User, homeTab, ev.View.Hash)
}
