package eventHandlers

import (
	"fmt"

	models "demoparser/models"
	utils "demoparser/utils"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

var DebugWeaponFire = false
var WeaponFirePlayerFilter = "" // If set, only log shots by this player

// HandleWeaponFireEvent processes WeaponFire events.
func HandleWeaponFireEvent(parser dem.Parser, event events.WeaponFire) {
	if event.Shooter != nil && event.Shooter.IsConnected {
		shooterID := event.Shooter.SteamID64
		currentTick := uint64(parser.GameState().IngameTick())

		// Log the shooter every time a weapon is fired if enabled and matches filter
		if DebugWeaponFire && (WeaponFirePlayerFilter == "" || event.Shooter.Name == WeaponFirePlayerFilter) {
			fmt.Printf("%s fired by: %s (SteamID: %d) at tick %d\n", event.Weapon.String(), event.Shooter.Name, shooterID, currentTick)
		}

		models.PlayerStatsMap[shooterID].WeaponFireCounts[event.Weapon.String()]++

		if utils.IsGrenade(event.Weapon) || utils.IsKnife(event.Weapon) {
			return // Ignore grenades and knives before incrementing shots
		}

		// Increment total shots fired
		models.PlayerStatsMap[shooterID].TotalShots++

		// Track the weapon fire for this tick
		if _, exists := models.WeaponFiredAtTick[shooterID]; !exists {
			models.WeaponFiredAtTick[shooterID] = make(map[uint64]bool)
		}
		models.WeaponFiredAtTick[shooterID][currentTick] = true

		// Ignore AWP shots
		if utils.IsAWP(event.Weapon) {
			return
		}

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
	}
}
