package karma

import (
	"fmt"
	"math/rand"
)

// Messages for different karma increments
var groupDoubleBoostMessages = []string{
	"The karma of <!subteam^%s> and its members got a double boost! 🚀 They now have %d karma.",
	"Someone must have turned up the cosmic amplifier because <!subteam^%s> and its members just received a double boost in karma! 🔊🌌 They now have %d karma.",
	"The cosmic karma dispenser malfunctioned and gave <!subteam^%s> and its members a double serving! 🤖🍲 They now have %d karma.",
	"It's raining karma for <!subteam^%s> and its members, and they just got caught in the downpour! ☔🌟 They now have %d karma.",
	"The universe must be feeling extra generous today because <!subteam^%s> and its members just hit the cosmic jackpot! 🎰🌠 They now have %d karma.",
	"Looks like the karma fairy made a special delivery to <!subteam^%s> and its members, and it was double the usual! 🧚‍♀️📦 They now have %d karma.",
	"The cosmic karma wave just swept through <!subteam^%s> and its members, and it was twice as big as usual! 🌊🚀 They now have %d karma.",
	"Did <!subteam^%s> and its members accidentally stumble upon a karma treasure chest? Because they just got double the loot! 🏴‍☠️💰 They now have %d karma.",
	"The karma accelerator just got stuck on overdrive for <!subteam^%s> and its members, doubling their speed to good fortune! 🏎️💨 They now have %d karma.",
	"The cosmic karma factory just had a production glitch, and <!subteam^%s> and its members received double the output! 🏭💫 They now have %d karma.",
}

var groupOnTheRiseMessages = []string{
	"The karma of <!subteam^%s> and its members is on the rise! 🚀 They now have %d karma.",
	"The karma barometer for <!subteam^%s> and its members is off the charts! 📈🌟 They now have %d karma.",
	"Looks like the cosmic karma stock for <!subteam^%s> and its members just skyrocketed! 📈🌌 They now have %d karma.",
	"The karma thermometer just broke because <!subteam^%s> and its members are heating up! 🌡️🔥 They now have %d karma.",
	"Did <!subteam^%s> and its members just turn on the karma engine? Because they're climbing fast! 🚗⬆️ They now have %d karma.",
	"The cosmic karma elevator just stopped at <!subteam^%s> and its members' floor! 🛗🌠 They now have %d karma.",
	"Looks like the karma escalator for <!subteam^%s> and its members just got an upgrade! 🛢️⬆️ They now have %d karma.",
	"The cosmic karma balloon for <!subteam^%s> and its members is inflating rapidly! 🎈🚀 They now have %d karma.",
	"It seems like <!subteam^%s> and its members just hit the jackpot in the karma lottery! 🎰✨ They now have %d karma.",
	"Looks like <!subteam^%s> and its members are riding the karma wave all the way to the top! 🌊🌟 They now have %d karma.",
}

var groupTookAHitMessages = []string{
	"The karma of <!subteam^%s> and its members took a hit! 💔 They now have %d karma.",
	"Uh-oh! It seems like the cosmic karma fairy just pranked <!subteam^%s> and its members! 🧚‍♀️😬 They now have %d karma.",
	"Looks like <!subteam^%s> and its members just stepped on a cosmic banana peel of bad luck! 🍌😬 They now have %d karma.",
	"The karma gods must be playing a joke on <!subteam^%s> and its members today! 😅🌌 They now have %d karma.",
	"The cosmic karma balance just got out of whack for <!subteam^%s> and its members! ⚖️🌟 They now have %d karma.",
	"Did <!subteam^%s> and its members accidentally break a mirror made of karma? Because they're experiencing double the bad luck! 🪞😬 They now have %d karma.",
	"Looks like <!subteam^%s> and its members just stumbled into the Bermuda Triangle of bad karma! 🌀😬 They now have %d karma.",
	"The cosmic karma magnet just attracted a storm of bad vibes to <!subteam^%s> and its members! 🧲⛈️ They now have %d karma.",
	"The universe just pulled a prank on <!subteam^%s> and its members by switching their karma with someone else's! 😅🌌 They now have %d karma.",
	"It seems like <!subteam^%s> and its members just entered the wrong door in the cosmic karma maze! 🚪😬 They now have %d karma.",
}

var groupDoubleHitMessages = []string{
	"The karma of <!subteam^%s> and its members took a double hit! 💔 They now have %d karma.",
	"Oh no! It looks like <!subteam^%s> and its members just got caught in a double whammy of bad karma! 😬💥 They now have %d karma.",
	"Did <!subteam^%s> and its members just step on a cosmic landmine of bad vibes? 💣😬 They now have %d karma.",
	"The cosmic karma storm just hit <!subteam^%s> and its members twice as hard! 🌪️💥 They now have %d karma.",
	"Looks like <!subteam^%s> and its members just walked through a double trouble field of cosmic karma! 👣😬 They now have %d karma.",
	"The universe just doubled down on bad luck for <!subteam^%s> and its members! 🎰😬 They now have %d karma.",
	"The cosmic karma avalanche just buried <!subteam^%s> and its members under a double layer of bad vibes! 🏔️😬 They now have %d karma.",
	"Uh-oh! It seems like <!subteam^%s> and its members just got caught in the crossfire of a double dose of cosmic karma! 🔥😬 They now have %d karma.",
	"The cosmic karma black hole just swallowed <!subteam^%s> and its members twice as deep! 🕳️😬 They now have %d karma.",
	"Looks like <!subteam^%s> and its members just stumbled into a double trouble vortex of bad karma! 🌀😬 They now have %d karma.",
}

// getGroupKarmaMessage returns a random message based on the karma increment
func getGroupKarmaMessage(groupId string, newKarma int, increment int) string {
	var messages []string

	switch increment {
	case 2:
		messages = groupDoubleBoostMessages
	case 1:
		messages = groupOnTheRiseMessages
	case -1:
		messages = groupTookAHitMessages
	case -2:
		messages = groupDoubleHitMessages
	}

	return fmt.Sprintf(messages[rand.Intn(len(messages))], groupId, newKarma)
}
