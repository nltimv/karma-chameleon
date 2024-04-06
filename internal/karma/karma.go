package karma

import (
	"fmt"
	"os"
	"regexp"

	"github.com/slack-go/slack/slackevents"
	"nltimv.com/karma-chameleon/internal/db"
	"nltimv.com/karma-chameleon/internal/slack"
)

func ProcessGetUserKarma(ev *slackevents.MessageEvent, apiEvent *slackevents.EventsAPIEvent, re *regexp.Regexp) {
	matches := re.FindStringSubmatch(ev.Text)
	userID := matches[1]

	karma := db.GetUserKarma(userID, apiEvent.TeamID)

	response := fmt.Sprintf("<@%s> currently has %d karma.", userID, karma)
	slack.Say(response, ev.Channel)
}

func ProcessGetGroupKarma(ev *slackevents.MessageEvent, apiEvent *slackevents.EventsAPIEvent, re *regexp.Regexp) {
	matches := re.FindStringSubmatch(ev.Text)
	groupID := matches[1]

	karma := db.GetGroupKarma(groupID, apiEvent.TeamID)

	response := fmt.Sprintf("<!subteam^%s> currently has %d karma.", groupID, karma)
	slack.Say(response, ev.Channel)
}

func ProcessUserKarma(ev *slackevents.MessageEvent, apiEvent *slackevents.EventsAPIEvent, re *regexp.Regexp) {
	matches := re.FindStringSubmatch(ev.Text)
	userID := matches[1]
	increment := matches[2]

	incrementValue := getIncrement(increment)

	if incrementValue == 0 {
		return
	}

	var karma int
	valid := slack.IsValidUser(userID)
	if valid {
		if userID != ev.User || incrementValue < 0 {
			var err error

			karma, err = db.UpdateUserKarma(userID, apiEvent.TeamID, incrementValue)
			if err != nil {
				fmt.Fprintf(os.Stdout, "Error while updating user karma: %v\n", err)
				return
			}
		} else {
			slack.Say("Nice try! You can't boost your own ego. 😜", ev.Channel)
			return
		}
	} else {
		fmt.Printf("Unknown user ID '%v'!\n", userID)
		return
	}

	response := getUserKarmaMessage(userID, karma, incrementValue)

	slack.Say(response, ev.Channel)
}

func ProcessGroupKarma(ev *slackevents.MessageEvent, apiEvent *slackevents.EventsAPIEvent, re *regexp.Regexp) {
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
		slack.Say("Nice try! Creating a user group for youself so you can get group karma? You're smart, but not smart enough! 😜", ev.Channel)
		return
	}

	var err error

	for _, memberID := range usergroupMembers {
		if memberID != ev.User || incrementValue < 0 {
			_, err = db.UpdateUserKarma(memberID, apiEvent.TeamID, incrementValue)
			if err != nil {
				fmt.Fprintf(os.Stdout, "Error while updating user karma: %v", err)
				return
			}
		}
	}

	karma := db.UpdateGroupKarma(groupID, apiEvent.TeamID, incrementValue)

	response := getGroupKarmaMessage(groupID, karma, incrementValue)

	slack.Say(response, ev.Channel)
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
