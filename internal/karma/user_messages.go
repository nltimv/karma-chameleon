package karma

import (
	"fmt"
	"math/rand"
)

// Messages for different karma increments
var userDoubleBoostMessages = []string{
	"<@%s>'s karma just unlocked the secret cheat code and got a double score! 🔓💥 They now have %d karma.",
	"<@%s>'s karma just did a victory dance and scored double points! 🎉💃 They now have %d karma.",
	"<@%s>'s karma just hit the jackpot and doubled up! 💰💰 They now have %d karma.",
	"<@%s>'s karma just got a double dose of good karma! ✨✨ They now have %d karma.",
	"<@%s>'s karma just found the genie lamp and wished for double karma! 🧞‍♂️✨ They now have %d karma.",
	"<@%s>'s karma just got upgraded to the deluxe edition, with double the benefits! 🎮💫 They now have %d karma.",
	"<@%s>'s karma just received a double high-five from the universe! 🖐️🌌 They now have %d karma.",
	"<@%s>'s karma just turned on the turbo boost and doubled its speed! 🚗💨 They now have %d karma.",
	"<@%s>'s karma just discovered the time warp and looped back for double rewards! ⏰🔄 They now have %d karma.",
	"<@%s>'s karma just got a two-for-one deal at the karma store! 🛍️🎁 They now have %d karma.",
}

var userOnTheRiseMessages = []string{
	"<@%s>'s karma is scaling the Everest of positivity! 🏔️ They now have %d karma.",
	"<@%s>'s karma is riding the wave of success! 🌊 They now have %d karma.",
	"<@%s>'s karma is soaring higher than a kite on a windy day! 🪁 They now have %d karma.",
	"<@%s>'s karma is climbing faster than Spider-Man up a skyscraper! 🕷️🏙️ They now have %d karma.",
	"<@%s>'s karma is shining brighter than a supernova! 💫 They now have %d karma.",
	"<@%s>'s karma is on the express train to Good-Vibes-Ville! 🚂💨 They now have %d karma.",
	"<@%s>'s karma is dancing its way to the top of the leaderboard! 💃🥇 They now have %d karma.",
	"<@%s>'s karma is cruising on the highway of happiness with the top down! 🛣️🌞 They now have %d karma.",
	"<@%s>'s karma is like a phoenix rising from the ashes of negativity! 🦅🔥 They now have %d karma.",
	"<@%s>'s karma is blossoming like a spring flower in full bloom! 🌸🌼 They now have %d karma.",
}

var userTookAHitMessages = []string{
	"<@%s>'s karma just hit a speed bump on the road of life! 🚧 They now have %d karma.",
	"<@%s>'s karma just got a detour through the valley of bad luck! ⛰️🍀 They now have %d karma.",
	"<@%s>'s karma just stepped on a LEGO in the dark! 😖 They now have %d karma.",
	"<@%s>'s karma just got caught in the crossfire of negativity! 💥 They now have %d karma.",
	"<@%s>'s karma just took a wrong turn at Albuquerque! 🗺️ They now have %d karma.",
	"<@%s>'s karma just got caught in a rainstorm without an umbrella! ☔ They now have %d karma.",
	"<@%s>'s karma just stubbed its toe on the corner of reality! 😣 They now have %d karma.",
	"<@%s>'s karma just hit a pothole on the highway to happiness! 🕳️ They now have %d karma.",
	"<@%s>'s karma just got tangled in the web of negativity! 🕸️ They now have %d karma.",
	"<@%s>'s karma just tripped over its own shoelaces! 👟 They now have %d karma.",
}

var userDoubleHitMessages = []string{
	"<@%s>'s karma just got sucker-punched twice by Murphy's Law! 👊💢 They now have %d karma.",
	"<@%s>'s karma just walked into a double dose of bad luck! 🍀🍀 They now have %d karma.",
	"<@%s>'s karma just stumbled into a double whammy of misfortune! 💥💔 They now have %d karma.",
	"<@%s>'s karma just got struck by lightning twice! ⚡⚡ They now have %d karma.",
	"<@%s>'s karma just got hit with a double scoop of negativity! 🍦🍦 They now have %d karma.",
	"<@%s>'s karma just took a double tumble down the stairs of fate! 🎢🎢 They now have %d karma.",
	"<@%s>'s karma just got hit by a double-decker bus of bad vibes! 🚌🚌 They now have %d karma.",
	"<@%s>'s karma just walked under two ladders and broke a mirror! 🪜🪜🪞 They now have %d karma.",
	"<@%s>'s karma just ran into a black cat and broke a mirror, twice! 🐈🪞🪞 They now have %d karma.",
	"<@%s>'s karma just stepped on two cracks and broke its mother's back, twice! 👩‍👦🪨🪨 They now have %d karma.",
}

// getUserKarmaMessage returns a random message depending on the karma increment
func getUserKarmaMessage(userId string, newKarma int, increment int) string {
	var messages []string

	switch increment {
	case 2:
		messages = userDoubleBoostMessages
	case 1:
		messages = userOnTheRiseMessages
	case -1:
		messages = userTookAHitMessages
	case -2:
		messages = userDoubleHitMessages
	}

	return fmt.Sprintf(messages[rand.Intn(len(messages))], userId, newKarma)
}
