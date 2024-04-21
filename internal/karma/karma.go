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

	user, err := db.GetUserKarma(userID, apiEvent.TeamID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while querying user karma: %v", err)
	}

	response := fmt.Sprintf("<@%s> currently has %d karma.", userID, user.Karma)
	slack.Say(response, ev.Channel)
}

func ProcessGetGroupKarma(ev *slackevents.MessageEvent, apiEvent *slackevents.EventsAPIEvent, re *regexp.Regexp) {
	matches := re.FindStringSubmatch(ev.Text)
	groupID := matches[1]

	group, err := db.GetGroupKarma(groupID, apiEvent.TeamID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while querying group karma: %v", err)
	}

	response := fmt.Sprintf("<!subteam^%s> currently has %d karma.", groupID, group.Karma)
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

	var user *db.User
	valid := slack.IsValidUser(userID)
	if valid {
		if userID != ev.User || incrementValue < 0 {
			var err error

			user, err = db.UpdateUserKarma(userID, apiEvent.TeamID, incrementValue)
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

	response := getUserKarmaMessage(userID, user.Karma, incrementValue)

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
				fmt.Fprintf(os.Stderr, "Error while updating user karma: %v", err)
				return
			}
		}
	}

	var group *db.Group

	if group, err = db.UpdateGroupKarma(groupID, apiEvent.TeamID, incrementValue); err != nil {
		fmt.Fprintf(os.Stderr, "Error while updating group karma: %v", err)
	}

	response := getGroupKarmaMessage(groupID, group.Karma, incrementValue)

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
