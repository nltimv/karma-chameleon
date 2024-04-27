package viewmodels

type LeaderboardUserRow struct {
	UserId            string
	Karma             int
	Rank              uint
	ProfilePictureUri string
}

type LeaderboardUserViewModel struct {
	CurrentUser *LeaderboardUserRow
	Leaderboard []*LeaderboardUserRow
}
