package db

import "gorm.io/gorm"

type UserLeaderboardEntry struct {
	User
	Rank uint
}

func GetUserLeaderboard(teamID string) ([]*UserLeaderboardEntry, error) {
	return handleTransaction(func(tx *gorm.DB) ([]*UserLeaderboardEntry, error) {
		var users []*UserLeaderboardEntry
		err := tx.Model(&User{}).Select("*, RANK () OVER (ORDER BY karma DESC) rank").Find(&users, &User{TeamId: teamID}).Error
		return users, err
	})
}
