package eventHandlers

import (
	"demoparser/models"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

// HandleKillEvent processes Kill events.
func HandleKillEvent(p dem.Parser, e events.Kill) {
	if e.Killer != nil && e.Killer.IsConnected {
		killerID := e.Killer.SteamID64
		if _, exists := models.PlayerStatsMap[killerID]; !exists {
			models.PlayerStatsMap[killerID] = &models.PlayerStats{SteamID: killerID, Name: e.Killer.Name}
		}
		models.PlayerStatsMap[killerID].Kills++
		if e.IsHeadshot {
			models.PlayerStatsMap[killerID].Headshots++
		}
	}
	if e.Victim != nil && e.Victim.IsConnected {
		victimID := e.Victim.SteamID64
		if _, exists := models.PlayerStatsMap[victimID]; !exists {
			models.PlayerStatsMap[victimID] = &models.PlayerStats{SteamID: victimID, Name: e.Victim.Name}
		}
		models.PlayerStatsMap[victimID].Deaths++
	}
}
