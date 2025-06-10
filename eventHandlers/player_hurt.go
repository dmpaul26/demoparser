package eventHandlers

import (
	"fmt"

	"demoparser/models"
	"demoparser/utils"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

// HandlePlayerHurtEvent processes PlayerHurt events.
func HandlePlayerHurtEvent(p dem.Parser, e events.PlayerHurt) {
	if e.Attacker != nil && e.Attacker.IsConnected && e.Player != nil && e.Player.IsConnected {
		attackerID := e.Attacker.SteamID64
		currentTick := uint64(p.GameState().IngameTick())

		if e.Attacker.Team == e.Player.Team {
			return // Ignore team damage
		}

		// Ignore grenades, knives
		if utils.IsGrenade(e.Weapon) || utils.IsKnife(e.Weapon) {
			return
		}

		if utils.IsShotgun(e.Weapon) {
			// Only count this shotgun shot once per tick
			if !models.ShotgunShots[currentTick] {
				models.ShotgunShots[currentTick] = true
				models.PlayerStatsMap[attackerID].ShotgunHits++
				models.PlayerStatsMap[attackerID].TotalHits++
			}
			if e.HitGroup == 1 && !models.ShotgunHSShots[currentTick] {
				models.ShotgunHSShots[currentTick] = true
				models.PlayerStatsMap[attackerID].ShotgunHSHits++
				models.PlayerStatsMap[attackerID].TotalHSHits++
			}
		} else {
			if !utils.IsAWP(e.Weapon) {
				// If it's NOT a shotgun, count normally
				models.PlayerStatsMap[attackerID].TotalHits++
				if e.HitGroup == 1 {
					models.PlayerStatsMap[attackerID].TotalHSHits++
				}
				if utils.IsRifle(e.Weapon) {
					models.PlayerStatsMap[attackerID].RifleHits++
					if e.HitGroup == 1 {
						models.PlayerStatsMap[attackerID].RifleHSHits++
					}
				}
			}
		}

		// Check if a weapon was fired within a range of ticks
		found := false
		for tickOffset := -3; tickOffset <= 3; tickOffset++ {
			tickToCheck := currentTick + uint64(tickOffset)
			if models.WeaponFiredAtTick[attackerID][tickToCheck] {
				if e.Attacker.Name == "itsPhix" {
					fmt.Printf("Weapon fire found for player %s at tick %d hitting player %s\n", e.Attacker.Name, tickToCheck, e.Player.Name)
				}
				delete(models.WeaponFiredAtTick[attackerID], tickToCheck)
				found = true
				break
			}
		}

		if !found {
			// Store the event for later processing
			models.PendingPlayerHurt[attackerID] = append(models.PendingPlayerHurt[attackerID], e)
			//fmt.Printf("No weapon fire found for player %s near tick %d\n", e.Attacker.Name, currentTick)
		}
	}
}
