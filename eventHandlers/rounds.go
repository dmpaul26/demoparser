package eventHandlers

import (
	"demoparser/models"
	"fmt"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

var DebugRounds = false
var DebugRoundEnds = false
var RoundCount int
var MatchStartCount int
var MaxRoundDebug = -1

// HandleRoundStartEvent logs the start of a round and increments the round count.
func HandleRoundStartEvent(parser dem.Parser, event events.RoundStart) {
	RoundCount++
	if DebugRounds {
		fmt.Printf("---------------------------- Round start: %d ----------------------------\n", RoundCount)
	}
	if MaxRoundDebug > 0 && RoundCount > MaxRoundDebug {
		toggleAllDebugs(false)
	}
}

func toggleAllDebugs(state bool) {
	DebugChat = state
	DebugRounds = state
	DebugRoundEnds = state
	DebugConnections = state
	DebugWeaponFire = state
	DebugSpotters = state
	models.DebugPlayerStatsInit = state // or false to disable
}

// HandleRoundEndEvent logs the end of a round and the current round count.
func HandleRoundEndEvent(parser dem.Parser, e events.RoundEnd) {
	if DebugRoundEnds {
		fmt.Printf("---------------------------- Round end: %d ----------------------------\n", RoundCount)
	}
}

// HandleMatchStartedEvent logs when the match starts and the current round count.
func HandleMatchStartedEvent(parser dem.Parser, e events.MatchStart) {
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
func HandleFinalRoundEvent(parser dem.Parser, e events.AnnouncementFinalRound) {
	if DebugRounds {
		fmt.Printf("Final round: %d\n", RoundCount)
	}
}

// HandleWinPanelMatchEvent logs when the win panel is shown and the current round count.
func HandleWinPanelMatchEvent(parser dem.Parser, e events.AnnouncementWinPanelMatch) {
	if DebugRounds {
		fmt.Printf("Win panel match: round %d\n", RoundCount)
	}
}
