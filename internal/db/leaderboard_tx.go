package db

import "gorm.io/gorm"

func GetUserLeaderboard(teamID string) ([]*UserLeaderboardEntry, error) {
	return handleTransaction(func(tx *gorm.DB) ([]*UserLeaderboardEntry, error) {
		return getUserLeaderboard(tx, teamID)
	})
}
