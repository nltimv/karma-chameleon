package slack

import (
	"nltimv.com/karma-chameleon/internal/log"
)

func IsValidUser(userID string) bool {
	userInfo, err := webApi.GetUserInfo(userID)
	if err != nil {
		log.Error.Println("Error getting user info: ", err)
		return false
	}

	return userInfo != nil && !userInfo.Deleted
}

func GetUsergroupMembers(groupID string) []string {
	usergroup, err := webApi.GetUserGroupMembers(groupID)
	if err != nil {
		log.Error.Println("Error getting user group members: ", err)
		return []string{}
	}

	return usergroup
}
