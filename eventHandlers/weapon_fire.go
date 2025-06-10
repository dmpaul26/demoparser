package eventHandlers

import (
	"fmt"

	models "demoparser/models"
	utils "demoparser/utils"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

// handleWeaponFireEvent processes WeaponFire events.
func HandleWeaponFireEvent(parser dem.Parser, event events.WeaponFire) {
	if event.Shooter != nil && event.Shooter.IsConnected {
		shooterID := event.Shooter.SteamID64
		currentTick := uint64(parser.GameState().IngameTick())

		// Initialize player stats if not already present
		if _, exists := models.PlayerStatsMap[shooterID]; !exists {
			models.PlayerStatsMap[shooterID] = &models.PlayerStats{SteamID: shooterID, Name: event.Shooter.Name}
		}
		// Ignore grenades, knives, and AWP shots
		if utils.IsGrenade(event.Weapon) || utils.IsKnife(event.Weapon) || utils.IsAWP(event.Weapon) {
			return
		}

		// Increment total shots fired
		models.PlayerStatsMap[shooterID].TotalShots++

		// Track the weapon fire for this tick
		if _, exists := models.WeaponFiredAtTick[shooterID]; !exists {
			models.WeaponFiredAtTick[shooterID] = make(map[uint64]bool)
		}
		models.WeaponFiredAtTick[shooterID][currentTick] = true

		// Process pending PlayerHurt events
		if events, exists := models.PendingPlayerHurt[shooterID]; exists {
			for _, hurtEvent := range events {
				hurtTick := uint64(parser.GameState().IngameTick())
				for tickOffset := -3; tickOffset <= 3; tickOffset++ {
					tickToCheck := hurtTick + uint64(tickOffset)
					if models.WeaponFiredAtTick[shooterID][tickToCheck] {
						if event.Shooter.Name == "itsPhix" {
							fmt.Printf("Weapon fire retro match with %s %s at tick %d hitting player %s\n", event.Weapon.String(), event.Shooter.Name, tickToCheck, hurtEvent.Player.Name)
						}
						delete(models.WeaponFiredAtTick[shooterID], tickToCheck)
						break
					}
				}
			}
			delete(models.PendingPlayerHurt, shooterID) // Clear processed events
		}

		// Log the weapon fire
		if event.Shooter.Name == "itsPhix" {
			fmt.Printf("%s fired by %s at tick %d\n", event.Weapon.String(), event.Shooter.Name, currentTick)
		}
	}
}
