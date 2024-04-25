package db

import "gorm.io/gorm"

func UpdateUserKarma(userID string, teamID string, increment int) (*User, error) {
	return handleTransaction(func(tx *gorm.DB) (*User, error) {
		return updateUserKarma(tx, userID, teamID, increment)
	})
}

func UpdateGroupKarma(groupID string, teamID string, increment int) (*Group, error) {
	return handleTransaction(func(tx *gorm.DB) (*Group, error) {
		return updateGroupKarma(tx, groupID, teamID, increment)
	})
}

func GetUserKarma(userID string, teamID string) (*User, error) {
	return handleTransaction(func(tx *gorm.DB) (*User, error) {
		return getUserKarma(tx, userID, teamID)
	})
}

func GetGroupKarma(groupID string, teamID string) (*Group, error) {
	return handleTransaction(func(tx *gorm.DB) (*Group, error) {
		return getGroupKarma(tx, groupID, teamID)
	})
}
