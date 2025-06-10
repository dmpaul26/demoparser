package eventHandlers

import (
	"fmt"

	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

var DebugRounds = false
var DebugRoundEnds = false
var RoundCount int
var MatchStartCount int

// HandleRoundStartEvent logs the start of a round and increments the round count.
func HandleRoundStartEvent(e events.RoundStart) {
	RoundCount++
	if DebugRounds {
		fmt.Printf("Round start: %d\n", RoundCount)
	}
}

// HandleRoundEndEvent logs the end of a round and the current round count.
func HandleRoundEndEvent(e events.RoundEnd) {
	if DebugRoundEnds {
		fmt.Printf("Round end: %d\n", RoundCount)
	}
}

// HandleMatchStartedEvent logs when the match starts and the current round count.
func HandleMatchStartedEvent(e events.MatchStart) {
	if DebugRounds {
		fmt.Printf("Match started: round %d\n", RoundCount)
	}

	RoundCount = 0
	MatchStartCount++

	if MatchStartCount > 1 {
		if DebugRounds {
			fmt.Printf("Warmup ended, Faceit match starting...\n")
		}
	}
}

// HandleFinalRoundEvent logs when the final round starts and the current round count.
func HandleFinalRoundEvent(e events.AnnouncementFinalRound) {
	if DebugRounds {
		fmt.Printf("Final round: %d\n", RoundCount)
	}
}

// HandleWinPanelMatchEvent logs when the win panel is shown and the current round count.
func HandleWinPanelMatchEvent(e events.AnnouncementWinPanelMatch) {
	if DebugRounds {
		fmt.Printf("Win panel match: round %d\n", RoundCount)
	}
}
