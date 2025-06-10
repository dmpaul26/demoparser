package models

import (
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

	// New fields for advanced stats
	TotalAimDistance float64 // Cumulative angular distance
	SpottingEvents   int     // Number of spotting events
	AverageAimDist   float64 // Average angular distance
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
