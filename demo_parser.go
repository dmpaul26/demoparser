package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

// PlayerStats struct to store kills, deaths, total HS acc, rifle HS acc, and shotgun HS hits
type PlayerStats struct {
	SteamID       uint64
	Name          string
	Kills         int
	Deaths        int
	Headshots     int
	RifleHSHits   int
	RifleHits     int
	TotalHSHits   int // HS hits from all weapons except grenades/knives/AWP
	TotalHits     int // All landed shots (excluding grenades/knives/AWP)
	ShotgunHits   int // Total grouped shotgun hits (one per tick)
	ShotgunHSHits int // Total grouped shotgun headshots (one per tick)
	HSPercentage  float64
	RifleHSAcc    float64
	TotalHSAcc    float64
	MissedShots   int // New field to track missed shots
	TotalShots    int // New field to track total shots fired
	TotalAcc      float64
}

// Player stats map (keyed by SteamID64 instead of name)
var playerStats = make(map[uint64]*PlayerStats)

// Track shotgun hits and headshots per tick
var shotgunShots = make(map[uint64]bool)
var shotgunHSShots = make(map[uint64]bool)

// Track weapon fires per tick
var weaponFiredAtTick = make(map[uint64]map[uint64]bool) // shooterID -> tick -> fired

var pendingPlayerHurt = make(map[uint64][]events.PlayerHurt) // attackerID -> list of PlayerHurt events

func isRifle(weapon *common.Equipment) bool {
	if weapon == nil {
		return false
	}
	rifles := map[string]bool{
		"AK-47":    true,
		"M4A4":     true,
		"M4A1":     true,
		"Galil AR": true,
		"FAMAS":    true,
	}
	return rifles[weapon.String()]
}

func isShotgun(weapon *common.Equipment) bool {
	if weapon == nil {
		return false
	}
	shotguns := map[string]bool{
		"XM1014":    true,
		"MAG-7":     true,
		"Nova":      true,
		"Sawed-Off": true,
	}
	return shotguns[weapon.String()]
}

func isAWP(weapon *common.Equipment) bool {
	return weapon != nil && weapon.String() == "AWP"
}

func isGrenade(weapon *common.Equipment) bool {
	if weapon == nil {
		return false
	}
	grenades := map[string]bool{
		"HE Grenade":         true,
		"Flashbang":          true,
		"Smoke Grenade":      true,
		"Decoy Grenade":      true,
		"Molotov":            true,
		"Incendiary Grenade": true,
	}
	return grenades[weapon.String()]
}

func isKnife(weapon *common.Equipment) bool {
	if weapon == nil {
		return false
	}
	knives := map[string]bool{
		"Knife":           true,
		"Bayonet":         true,
		"M9 Bayonet":      true,
		"Karambit":        true,
		"Butterfly Knife": true,
		"Flip Knife":      true,
	}
	return knives[weapon.String()]
}

func processDemo(demoPath string) {
	f, err := os.Open(demoPath)
	if err != nil {
		log.Panicf("Failed to open demo file %s: %v\n", demoPath, err)
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()

	fmt.Printf("Processing demo: %s\n", demoPath)

	// Register handler for kill events
	p.RegisterEventHandler(func(e events.Kill) {
		if e.Killer != nil && e.Killer.IsConnected {
			killerID := e.Killer.SteamID64
			if _, exists := playerStats[killerID]; !exists {
				playerStats[killerID] = &PlayerStats{SteamID: killerID, Name: e.Killer.Name}
			}
			playerStats[killerID].Kills++
			if e.IsHeadshot {
				playerStats[killerID].Headshots++
			}
		}
		if e.Victim != nil && e.Victim.IsConnected {
			victimID := e.Victim.SteamID64
			if _, exists := playerStats[victimID]; !exists {
				playerStats[victimID] = &PlayerStats{SteamID: victimID, Name: e.Victim.Name}
			}
			playerStats[victimID].Deaths++
		}
	})

	// Register handler for player hurt events
	p.RegisterEventHandler(func(e events.PlayerHurt) {
		if e.Attacker != nil && e.Attacker.IsConnected && e.Player != nil && e.Player.IsConnected {
			if e.Attacker.Team == e.Player.Team {
				return // Ignore team damage
			}

			attackerID := e.Attacker.SteamID64
			currentTick := uint64(p.GameState().IngameTick())

			if _, exists := playerStats[attackerID]; !exists {
				playerStats[attackerID] = &PlayerStats{SteamID: attackerID, Name: e.Attacker.Name}
			}

			// Ignore grenades, knives, and AWP shots
			if isGrenade(e.Weapon) || isKnife(e.Weapon) {
				return
			}

			if e.Attacker.Name == "itsPhix" {
				fmt.Printf("PlayerHurt event: Attacker %s, Victim %s, Tick %d with %s\n", e.Attacker.Name, e.Player.Name, currentTick, e.Weapon.String())
			}

			if isShotgun(e.Weapon) {
				// Only count this shotgun shot once per tick
				if !shotgunShots[currentTick] {
					shotgunShots[currentTick] = true
					playerStats[attackerID].ShotgunHits++
					playerStats[attackerID].TotalHits++ // Now counting only if it's a new shotgun shot
				}
				if e.HitGroup == 1 && !shotgunHSShots[currentTick] {
					shotgunHSShots[currentTick] = true
					playerStats[attackerID].ShotgunHSHits++
					playerStats[attackerID].TotalHSHits++
				}
			} else {
				if !isAWP(e.Weapon) {
					// If it's NOT a shotgun, count normally
					playerStats[attackerID].TotalHits++
					if e.HitGroup == 1 {
						playerStats[attackerID].TotalHSHits++
					}
					if isRifle(e.Weapon) {
						playerStats[attackerID].RifleHits++
						if e.HitGroup == 1 {
							playerStats[attackerID].RifleHSHits++
						}
					}
				}
			}

			// Check if a weapon was fired within a range of ticks
			found := false
			for tickOffset := -3; tickOffset <= 3; tickOffset++ {
				tickToCheck := currentTick + uint64(tickOffset)
				if weaponFiredAtTick[attackerID][tickToCheck] {
					if e.Attacker.Name == "itsPhix" {
						fmt.Printf("Weapon fire found for player %s at tick %d hitting player %s\n", e.Attacker.Name, tickToCheck, e.Player.Name)
					}
					delete(weaponFiredAtTick[attackerID], tickToCheck)
					found = true
					break
				}
			}

			if !found {
				// Store the event for later processing
				pendingPlayerHurt[attackerID] = append(pendingPlayerHurt[attackerID], e)
				//fmt.Printf("No weapon fire found for player %s near tick %d\n", e.Attacker.Name, currentTick)
			}
		}
	})

	// Register handler for weapon fire events
	p.RegisterEventHandler(func(e events.WeaponFire) {
		if e.Shooter != nil && e.Shooter.IsConnected {
			shooterID := e.Shooter.SteamID64
			currentTick := uint64(p.GameState().IngameTick())

			// Initialize player stats if not already present
			if _, exists := playerStats[shooterID]; !exists {
				playerStats[shooterID] = &PlayerStats{SteamID: shooterID, Name: e.Shooter.Name}
			}
			// Ignore grenades, knives, and AWP shots
			if isGrenade(e.Weapon) || isKnife(e.Weapon) || isAWP(e.Weapon) {
				return
			}

			// Increment total shots fired
			playerStats[shooterID].TotalShots++

			// Track the weapon fire for this tick
			if _, exists := weaponFiredAtTick[shooterID]; !exists {
				weaponFiredAtTick[shooterID] = make(map[uint64]bool)
			}
			weaponFiredAtTick[shooterID][currentTick] = true

			// Process pending PlayerHurt events
			if events, exists := pendingPlayerHurt[shooterID]; exists {
				for _, hurtEvent := range events {
					hurtTick := uint64(p.GameState().IngameTick())
					for tickOffset := -3; tickOffset <= 3; tickOffset++ {
						tickToCheck := hurtTick + uint64(tickOffset)
						if weaponFiredAtTick[shooterID][tickToCheck] {
							if e.Shooter.Name == "itsPhix" {
								fmt.Printf("Weapon fire retro match with %s %s at tick %d hitting player %s\n", e.Weapon.String(), e.Shooter.Name, tickToCheck, hurtEvent.Player.Name)
							}
							delete(weaponFiredAtTick[shooterID], tickToCheck)
							break
						}
					}
				}
				delete(pendingPlayerHurt, shooterID) // Clear processed events
			}

			// Log the weapon fire
			if e.Shooter.Name == "itsPhix" {
				fmt.Printf("%s fired by %s at tick %d\n", e.Weapon.String(), e.Shooter.Name, currentTick)
			}
		}
	})

	err = p.ParseToEnd()
	if err != nil {
		log.Panicf("Failed to parse demo %s: %v\n", demoPath, err)
	}
}

func main() {
	files, err := filepath.Glob("*.dem")
	if err != nil {
		log.Panic("Failed to read demo files: ", err)
	}

	if len(files) == 0 {
		fmt.Println("No .dem files found in the current folder.")
		return
	}

	for _, demoFile := range files {
		processDemo(demoFile)
	}

	var sortedPlayers []*PlayerStats
	for _, stats := range playerStats {
		sortedPlayers = append(sortedPlayers, stats)
	}

	// Count missed shots for each player
	for shooterID, ticks := range weaponFiredAtTick {
		if stats, exists := playerStats[shooterID]; exists {
			stats.MissedShots += len(ticks) // Each remaining tick represents a missed shot
		}
	}

	for _, stats := range playerStats {
		// Headshot percentage based on total kills
		if stats.Kills > 0 {
			stats.HSPercentage = (float64(stats.Headshots) / float64(stats.Kills)) * 100
		} else {
			stats.HSPercentage = 0
		}

		// Rifle Headshot Accuracy
		if stats.RifleHits > 0 {
			stats.RifleHSAcc = (float64(stats.RifleHSHits) / float64(stats.RifleHits)) * 100
		} else {
			stats.RifleHSAcc = 0
		}

		// Total Headshot Accuracy (excluding grenades, knives, AWP)
		if stats.TotalHits > 0 {
			stats.TotalHSAcc = (float64(stats.TotalHSHits) / float64(stats.TotalHits)) * 100
		} else {
			stats.TotalHSAcc = 0
		}

		// Total Accuracy
		if stats.TotalShots > 0 {
			stats.TotalAcc = ((float64(stats.TotalShots-stats.MissedShots) / float64(stats.TotalShots)) * 100)
		} else {
			stats.TotalAcc = 0
		}
	}

	sort.Slice(sortedPlayers, func(i, j int) bool {
		return sortedPlayers[i].TotalHSAcc > sortedPlayers[j].TotalHSAcc
	})

	fmt.Println("\nSorted Player Stats (by Total Headshot Acc %):")
	fmt.Printf("%-20s %-8s %-8s %-10s %-15s %-15s %-10s\n", "Player", "Kills", "Deaths", "HS%", "Rifle HS Acc", "Total HS Acc", "Accuracy")
	fmt.Println(strings.Repeat("-", 100))
	for _, stats := range sortedPlayers {
		fmt.Printf("%-20s %-8d %-8d %-10.2f %-15.2f %-15.2f %-10.2f\n",
			stats.Name, stats.Kills, stats.Deaths, stats.HSPercentage, stats.RifleHSAcc, stats.TotalHSAcc, stats.TotalAcc)
	}
}
