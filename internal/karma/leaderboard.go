package karma

import (
	slackapi "github.com/slack-go/slack"
	"nltimv.com/karma-chameleon/internal/db"
	"nltimv.com/karma-chameleon/internal/slack"
	"nltimv.com/karma-chameleon/internal/ui"
	"nltimv.com/karma-chameleon/internal/ui/viewmodels"
)

func OpenLeaderboard(interaction *slackapi.InteractionCallback) {
	// Get leaderboard entries from the database
	entries, err := db.GetUserLeaderboard(interaction.Team.ID)

	if err != nil {
		// Handle error
		return
	}

	// Create leaderboard view model
	leaderboard := make([]*viewmodels.LeaderboardUserRow, len(entries))

	var currentUserEntry *viewmodels.LeaderboardUserRow

	for i, entry := range entries {
		// Fill leaderboard view model
		leaderboard[i] = &viewmodels.LeaderboardUserRow{
			Rank:   entry.Rank,
			UserId: entry.UserId,
			Karma:  entry.Karma,
		}

		// Retrieve profile picture URI for top 3 entries
		if entry.Rank <= 3 {
			profilePictureURI, err := slack.GetProfilePictureUri(entry.UserId)

			if err != nil {
				// Handle error
				return
			}

			leaderboard[i].ProfilePictureUri = profilePictureURI
		}

		if entry.UserId == interaction.User.ID {
			currentUserEntry = leaderboard[i]
		}
	}

	// Create default leaderboard entry for the current user if not found
	if currentUserEntry == nil {
		currentUserEntry = &viewmodels.LeaderboardUserRow{
			UserId: interaction.User.ID,
			Karma:  0,
			Rank:   0,
		}
	}

	// Create leaderboard view model
	viewModel := &viewmodels.LeaderboardUserViewModel{
		CurrentUser: currentUserEntry,
		Leaderboard: leaderboard,
	}

	// Show leaderboard
	ui.ShowUserLeaderboard(viewModel, interaction)
}
