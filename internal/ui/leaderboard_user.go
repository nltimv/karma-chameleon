package ui

import (
	"fmt"

	slackapi "github.com/slack-go/slack"
	"nltimv.com/karma-chameleon/internal/slack"
	"nltimv.com/karma-chameleon/internal/ui/viewmodels"
)

func ShowUserLeaderboard(viewmodel *viewmodels.LeaderboardUserViewModel, interaction *slackapi.InteractionCallback) {
	blocks := make([]slackapi.Block, 0)

	blocks = append(blocks, slackapi.NewHeaderBlock(
		slackapi.NewTextBlockObject("plain_text", ":star: Your karma", true, false),
	))

	blocks = append(blocks, slackapi.NewSectionBlock(
		slackapi.NewTextBlockObject("mrkdwn", fmt.Sprintf("You are currently #%d on the leaderboard with %d karma.", viewmodel.CurrentUser.Rank, viewmodel.CurrentUser.Karma), false, false),
		nil, nil,
	))

	blocks = append(blocks, slackapi.NewDividerBlock())

	blocks = append(blocks, slackapi.NewHeaderBlock(
		slackapi.NewTextBlockObject("plain_text", ":trophy: Leaderboard", true, false),
	))

	for _, entry := range viewmodel.Leaderboard {
		var entrySection *slackapi.SectionBlock
		if entry.Rank <= 3 {
			text := fmt.Sprintf("*#%d* %v\n\t<@%s>\n\t%d karma", entry.Rank, getMedalEmoji(entry.Rank), entry.UserId, entry.Karma)
			textBlock := slackapi.NewTextBlockObject("mrkdwn", text, false, false)
			pfpBlock := slackapi.NewImageBlockElement(entry.ProfilePictureUri, "profile picture")
			entrySection = slackapi.NewSectionBlock(textBlock, nil, slackapi.NewAccessory(pfpBlock))
		} else {
			leftField := slackapi.NewTextBlockObject("mrkdwn", fmt.Sprintf("*#%d*\t<@%s>", entry.Rank, entry.UserId), false, false)
			rightField := slackapi.NewTextBlockObject("mrkdwn", fmt.Sprintf("%d karma", entry.Karma), false, false)
			fields := [2]*slackapi.TextBlockObject{leftField, rightField}
			entrySection = slackapi.NewSectionBlock(nil, fields[:], nil)
		}

		blocks = append(blocks, entrySection)
	}

	close := slackapi.NewTextBlockObject("plain_text", "Close", true, false)
	title := slackapi.NewTextBlockObject("plain_text", "Leaderboard (users)", true, false)

	mvr := slackapi.ModalViewRequest{
		Type:   slackapi.VTModal,
		Blocks: slackapi.Blocks{BlockSet: blocks},
		Close:  close,
		Title:  title,
	}

	slack.OpenModal(interaction.TriggerID, mvr)
}

func getMedalEmoji(rank uint) string {
	switch rank {
	case 1:
		return ":first_place_medal:"
	case 2:
		return ":second_place_medal:"
	case 3:
		return ":third_place_medal:"
	default:
		return ""
	}
}
