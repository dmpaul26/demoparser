package eventHandlers

import (
	"fmt"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

var DebugSpotters = false

// HandlePlayerSpottersChangedEvent processes PlayerSpottersChanged events.
func HandlePlayerSpottersChangedEvent(parser dem.Parser, event events.PlayerSpottersChanged) {
	if DebugSpotters {
		spotters := parser.GameState().Participants().SpottersOf(event.Spotted)
		if len(spotters) == 0 {
			fmt.Printf("Spotters disappeared: %s at tick %d\n", event.Spotted, parser.GameState().IngameTick())
		} else {
			fmt.Printf("Spotters changed: %s at tick %d\n", event.Spotted, parser.GameState().IngameTick())
			fmt.Printf("%s's spotters:\n", event.Spotted)
			for _, spotter := range spotters {
				if spotter != nil {
					fmt.Println(spotter.Name)
				}
			}
		}
	}

	// this is all broke idk yet
	/* if e.Spotted == nil || !e.Spotted.IsConnected || !e.Spotted.IsAlive() {
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
		if !e.Spotted.IsConnected || !player.IsAlive() || player.PlayerPawnEntity() == nil || e.Spotted == nil {
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

			// Update stats
			models.PlayerStatsMap[spotterID].TotalAimDistance += angularDistance
			models.PlayerStatsMap[spotterID].SpottingEvents++
		}
	}

	// Update the spotters map for this player
	models.PlayerSpotters[spottedPlayerID] = currentSpotters */
}
