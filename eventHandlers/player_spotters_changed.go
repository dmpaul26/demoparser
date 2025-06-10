package eventHandlers

import (
	"math"

	"demoparser/models"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

// HandlePlayerSpottersChangedEvent processes PlayerSpottersChanged events.
func HandlePlayerSpottersChangedEvent(p dem.Parser, e events.PlayerSpottersChanged) {
	if e.Spotted == nil || !e.Spotted.IsConnected || !e.Spotted.IsAlive() {
		return
	}

	spottedPlayerID := e.Spotted.SteamID64
	spottedPlayerPos := e.Spotted.Position()

	// Initialize the spotters map for this player if not already present
	if _, exists := models.PlayerSpotters[spottedPlayerID]; !exists {
		models.PlayerSpotters[spottedPlayerID] = make(map[uint64]bool)
	}

	// Get the current spotters from the game state
	currentSpotters := make(map[uint64]bool)
	for _, player := range p.GameState().Participants().Playing() {
		if player == nil || !player.IsConnected || !e.Spotted.IsConnected || !player.IsAlive() || player.PlayerPawnEntity() == nil || e.Spotted == nil {
			continue
		}

		// Check if the player is spotting the spotted player
		if player.IsSpottedBy(e.Spotted) {
			currentSpotters[player.SteamID64] = true
		}
	}

	// Compare current spotters with previously tracked spotters
	for spotterID := range currentSpotters {
		if !models.PlayerSpotters[spottedPlayerID][spotterID] {
			// New spotter detected, calculate angular distance
			spotter := p.GameState().Participants().FindByHandle64(spotterID)
			if spotter == nil || !spotter.IsConnected || !spotter.IsAlive() {
				continue
			}

			spotterPos := spotter.Position()
			spotterViewAngles := spotter.ViewDirectionX()

			deltaX := spottedPlayerPos.X - spotterPos.X
			deltaY := spottedPlayerPos.Y - spotterPos.Y
			angleToEnemy := math.Atan2(deltaY, deltaX) * (180 / math.Pi) // Convert to degrees
			angularDistance := math.Abs(angleToEnemy - float64(spotterViewAngles))

			// Normalize the angular distance to [0, 180]
			if angularDistance > 180 {
				angularDistance = 360 - angularDistance
			}

			// Initialize player stats if not already present
			if _, exists := models.PlayerStatsMap[spotterID]; !exists {
				models.PlayerStatsMap[spotterID] = &models.PlayerStats{SteamID: spotterID, Name: spotter.Name}
			}

			// Update stats
			models.PlayerStatsMap[spotterID].TotalAimDistance += angularDistance
			models.PlayerStatsMap[spotterID].SpottingEvents++
		}
	}

	// Update the spotters map for this player
	models.PlayerSpotters[spottedPlayerID] = currentSpotters
}
