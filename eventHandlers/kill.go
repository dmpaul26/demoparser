package eventHandlers

import (
	"demoparser/models"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

// HandleKillEvent processes Kill events.
func HandleKillEvent(parser dem.Parser, event events.Kill) {
	if event.Killer != nil && event.Killer.IsConnected {
		killerID := event.Killer.SteamID64
		models.PlayerStatsMap[killerID].Kills++
		if event.IsHeadshot {
			models.PlayerStatsMap[killerID].Headshots++
		}
	}
	if event.Victim != nil && event.Victim.IsConnected {
		victimID := event.Victim.SteamID64
		models.PlayerStatsMap[victimID].Deaths++
	}
}
