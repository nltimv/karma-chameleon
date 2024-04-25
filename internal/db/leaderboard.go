package db

import "gorm.io/gorm"

type UserLeaderboardEntry struct {
	User
	Rank int
}

func getUserLeaderboard(tx *gorm.DB, teamID string) ([]*UserLeaderboardEntry, error) {
	var users []*UserLeaderboardEntry
	err := tx.Model(&User{}).Select("*, RANK () OVER (ORDER BY karma DESC) rank").Find(&users, &User{TeamId: teamID}).Error
	return users, err
}
