package karma

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/slack-go/slack/slackevents"
	"nltimv.com/karma-chameleon/internal/db"
	"nltimv.com/karma-chameleon/internal/log"
	"nltimv.com/karma-chameleon/internal/slack"
)

func ProcessGetUserKarma(ev *slackevents.MessageEvent, apiEvent *slackevents.EventsAPIEvent, re *regexp.Regexp) {
	matches := re.FindStringSubmatch(ev.Text)
	userIDs := matches[1:]

	for _, userID := range userIDs {
		user, err := db.GetUserKarma(userID, apiEvent.TeamID)
		if err != nil {
			log.Error.Printf("Error while querying user karma: %v\n", err)
		}

		response := fmt.Sprintf("<@%s> currently has %d karma.", userID, user.Karma)
		ctx := getMessageContext(ev, apiEvent)
		slack.Say(response, ctx, nil)
	}
}

func ProcessGetGroupKarma(ev *slackevents.MessageEvent, apiEvent *slackevents.EventsAPIEvent, re *regexp.Regexp) {
	matches := re.FindStringSubmatch(ev.Text)
	groupID := matches[1]

	group, err := db.GetGroupKarma(groupID, apiEvent.TeamID)
	if err != nil {
		log.Error.Printf("Error while querying group karma: %v\n", err)
	}

	response := fmt.Sprintf("<!subteam^%s> currently has %d karma.", groupID, group.Karma)
	ctx := getMessageContext(ev, apiEvent)
	slack.Say(response, ctx, nil)
}

func ProcessUserKarma(ev *slackevents.MessageEvent, apiEvent *slackevents.EventsAPIEvent, re *regexp.Regexp) {
	ctx := getMessageContext(ev, apiEvent)
	matches := re.FindStringSubmatch(ev.Text)
	userRe := regexp.MustCompile(`<@([a-zA-Z0-9_]+)>`)
	userIDs := userRe.FindAllStringSubmatch(matches[0], -1)
	increment := matches[len(matches)-1]

	incrementValue := getIncrement(increment)

	if incrementValue == 0 {
		return
	}

	responses := make([]string, len(userIDs))

	for i, u := range userIDs {
		var user *db.User
		userID := u[1]
		valid := slack.IsValidUser(userID)
		if valid {
			if userID != ev.User || incrementValue < 0 {
				var err error

				user, err = db.UpdateUserKarma(userID, apiEvent.TeamID, incrementValue)
				if err != nil {
					log.Error.Printf("Error while updating user karma: %v\n", err)
					return
				}
			} else {
				slack.Say("Nice try! You can't boost your own ego. 😜", ctx, nil)
				return
			}
		} else {
			log.Error.Printf("Unknown user ID '%v'!\n", userID)
			return
		}

		switch incrementValue {
		case 2:
			responses[i] = fmt.Sprintf("<@%s>'s karma got a double boost! 🚀 They now have %d karma.", userID, user.Karma)
		case 1:
			responses[i] = fmt.Sprintf("<@%s>'s karma is on the rise! 🚀 They now have %d karma.", userID, user.Karma)
		case -1:
			responses[i] = fmt.Sprintf("<@%s>'s karma took a hit! 💔 They now have %d karma.", userID, user.Karma)
		case -2:
			responses[i] = fmt.Sprintf("<@%s>'s karma took a double hit! 💔 They now have %d karma.", userID, user.Karma)
		}
	}

	response := strings.Join(responses, "\n")

	slack.Say(response, ctx, nil)
}

func ProcessGroupKarma(ev *slackevents.MessageEvent, apiEvent *slackevents.EventsAPIEvent, re *regexp.Regexp) {
	ctx := getMessageContext(ev, apiEvent)
	matches := re.FindStringSubmatch(ev.Text)
	groupID := matches[1]
	increment := matches[2]

	incrementValue := getIncrement(increment)

	if incrementValue == 0 {
		return
	}

	usergroupMembers := slack.GetUsergroupMembers(groupID)

	if len(usergroupMembers) == 0 {
		return
	} else if len(usergroupMembers) == 1 && usergroupMembers[0] == ev.User {
		slack.Say("Nice try! Creating a user group for youself so you can get group karma? You're smart, but not smart enough! 😜", ctx, nil)
		return
	}

	var err error

	for _, memberID := range usergroupMembers {
		if memberID != ev.User || incrementValue < 0 {
			_, err = db.UpdateUserKarma(memberID, apiEvent.TeamID, incrementValue)
			if err != nil {
				log.Error.Printf("Error while updating user karma: %v\n", err)
				return
			}
		}
	}

	var group *db.Group

	if group, err = db.UpdateGroupKarma(groupID, apiEvent.TeamID, incrementValue); err != nil {
		log.Error.Printf("Error while updating group karma: %v", err)
	}

	var response string
	switch incrementValue {
	case 2:
		response = fmt.Sprintf("The karma of <!subteam^%s> and its members got a double boost! 🚀 They now have %d karma.", groupID, group.Karma)
	case 1:
		response = fmt.Sprintf("The karma of <!subteam^%s> and its members is on the rise! 🚀 They now have %d karma.", groupID, group.Karma)
	case -1:
		response = fmt.Sprintf("The karma of <!subteam^%s> and its members took a hit! 💔 They now have %d karma.", groupID, group.Karma)
	case -2:
		response = fmt.Sprintf("The karma of <!subteam^%s> and its members took a double hit! 💔 They now have %d karma.", groupID, group.Karma)
	}

	slack.Say(response, ctx, nil)
}

func getIncrement(incrString string) int {
	var incrementValue int
	switch incrString {
	case "+++":
		incrementValue = 2
	case "++":
		incrementValue = 1
	case "---":
		incrementValue = -2
	case "--":
		incrementValue = -1
	default:
		incrementValue = 0
	}

	return incrementValue
}

func getMessageContext(ev *slackevents.MessageEvent, apiEvent *slackevents.EventsAPIEvent) *slack.MessageContext {
	return &slack.MessageContext{
		ChannelId: ev.Channel,
		TeamId:    apiEvent.TeamID,
		AppId:     apiEvent.APIAppID,
	}
}
