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
	leaderboard := make([]*viewmodels.LeaderboardUserRow, 0)

	var currentUserEntry *viewmodels.LeaderboardUserRow

	for _, entry := range entries {
		// Fill leaderboard view model
		modelEntry := &viewmodels.LeaderboardUserRow{
			Rank:   entry.Rank,
			UserId: entry.UserId,
			Karma:  entry.Karma,
		}

		// Retrieve profile picture URI for top 3 entries
		if entry.Rank <= 3 || entry.UserId == interaction.User.ID {
			profilePictureURI, err := slack.GetProfilePictureUri(entry.UserId)

			if err != nil {
				// Handle error
				return
			}

			modelEntry.ProfilePictureUri = profilePictureURI
		}

		if entry.UserId == interaction.User.ID {
			currentUserEntry = modelEntry
		}

		if entry.Karma > 0 {
			leaderboard = append(leaderboard, modelEntry)
		}
	}

	// Create default leaderboard entry for the current user if not found
	if currentUserEntry == nil {
		profilePictureURI, err := slack.GetProfilePictureUri(interaction.User.ID)

		if err != nil {
			// Handle error
			return
		}

		currentUserEntry = &viewmodels.LeaderboardUserRow{
			UserId:            interaction.User.ID,
			Karma:             0,
			Rank:              0,
			ProfilePictureUri: profilePictureURI,
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
