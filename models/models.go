package models

import (
	"fmt"

	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

// DebugPlayerStatsInit controls logging for player stats initialization.
var DebugPlayerStatsInit = false

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

	// New fields for advanced stats
	TotalAimDistance float64 // Cumulative angular distance
	SpottingEvents   int     // Number of spotting events
	AverageAimDist   float64 // Average angular distance

	// WeaponFireCounts tracks the number of times each weapon has been fired
	WeaponFireCounts map[string]int
}

// Player stats map (keyed by SteamID64 instead of name)
var PlayerStatsMap = make(map[uint64]*PlayerStats)

// Track shotgun hits and headshots per tick
var ShotgunShots = make(map[uint64]bool)
var ShotgunHSShots = make(map[uint64]bool)

// Track weapon fires per tick
var WeaponFiredAtTick = make(map[uint64]map[uint64]bool) // shooterID -> tick -> fired

var PendingPlayerHurt = make(map[uint64][]events.PlayerHurt) // attackerID -> list of PlayerHurt events

var PlayerSpotters = make(map[uint64]map[uint64]bool) // playerID -> set of spotterIDs

func TryInitializePlayerStatsMap(shooterID uint64, name string) {
	if stats, exists := PlayerStatsMap[shooterID]; !exists {
		if DebugPlayerStatsInit {
			fmt.Printf("Initializing PlayerStats for %s (SteamID: %d)\n", name, shooterID)
		}
		PlayerStatsMap[shooterID] = &PlayerStats{
			SteamID:          shooterID,
			Name:             name,
			WeaponFireCounts: make(map[string]int),
		}
	} else if stats.Name != name {
		if DebugPlayerStatsInit {
			fmt.Printf("Updating PlayerStats name for SteamID %d: '%s' -> '%s'\n", shooterID, stats.Name, name)
		}
		stats.Name = name
	}
}
