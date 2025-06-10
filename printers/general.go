package printers

import (
	"demoparser/models"
	"fmt"
	"sort"
	"strings"
)

func PrintStats() {
	var sortedPlayers []*models.PlayerStats
	for _, stats := range models.PlayerStatsMap {
		sortedPlayers = append(sortedPlayers, stats)
	}

	// Count missed shots for each player
	for shooterID, ticks := range models.WeaponFiredAtTick {
		if stats, exists := models.PlayerStatsMap[shooterID]; exists {
			stats.MissedShots += len(ticks) // Each remaining tick represents a missed shot
		}
	}

	for _, stats := range models.PlayerStatsMap {
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

		// Calculate average angular distance
		if stats.SpottingEvents > 0 {
			stats.AverageAimDist = stats.TotalAimDistance / float64(stats.SpottingEvents)
		} else {
			stats.AverageAimDist = 0
		}
	}

	sort.Slice(sortedPlayers, func(i, j int) bool {
		return sortedPlayers[i].TotalHSAcc > sortedPlayers[j].TotalHSAcc
	})

	fmt.Println("\nSorted Player Stats (by Total Headshot Acc %):")
	fmt.Printf("%-20s %-8s %-8s %-10s %-15s %-15s %-10s %-15s\n", "Player", "Kills", "Deaths", "HS%", "Rifle HS Acc", "Total HS Acc", "Accuracy", "Avg Aim Dist")
	fmt.Println(strings.Repeat("-", 120))
	for _, stats := range sortedPlayers {
		fmt.Printf("%-20s %-8d %-8d %-10.2f %-15.2f %-15.2f %-10.2f %-15.2f\n",
			stats.Name, stats.Kills, stats.Deaths, stats.HSPercentage, stats.RifleHSAcc, stats.TotalHSAcc, stats.TotalAcc, stats.AverageAimDist)
	}
}
