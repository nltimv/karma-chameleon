package slack

import (
	"log"

	"github.com/slack-go/slack"
)

func Say(message string, channel string) {
	_, _, err := webApi.PostMessage(channel, slack.MsgOptionText(message, false))
	if err != nil {
		log.Printf("failed posting message: %v\n", err)
	}
}

func IsValidUser(userID string) bool {
	userInfo, err := webApi.GetUserInfo(userID)
	if err != nil {
		log.Println("Error getting user info: ", err)
		return false
	}

	return userInfo != nil && !userInfo.Deleted
}

func GetUsergroupMembers(groupID string) []string {
	usergroup, err := webApi.GetUserGroupMembers(groupID)
	if err != nil {
		log.Println("Error getting user group members: ", err)
		return []string{}
	}

	return usergroup
}
