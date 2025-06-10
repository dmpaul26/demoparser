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

// weaponCategories defines the order and grouping for weapon types.
var weaponCategories = [][]string{
	// Pistols
	{"Glock-18", "P2000", "USP-S", "P250", "Five-SeveN", "CZ75 Auto", "Tec-9", "Dual Berettas", "Desert Eagle", "R8 Revolver"},
	// Heavy
	{"Nova", "XM1014", "MAG-7", "Sawed-Off", "M249", "Negev"},
	// SMG
	{"MAC-10", "MP9", "MP7", "MP5-SD", "UMP-45", "P90", "PP-Bizon"},
	// Rifle
	{"FAMAS", "Galil AR", "M4A4", "M4A1-S", "M4A1", "AK-47", "AUG", "SG 553"},
	// Sniper
	{"SSG 08", "AWP", "G3SG1", "SCAR-20"},
	// Utility
	{"HE Grenade", "Flashbang", "Smoke Grenade", "Molotov", "Incendiary Grenade", "Decoy Grenade"},
	// Knives (all types)
	{
		"Knife", "Bayonet", "Flip Knife", "Gut Knife", "Karambit", "M9 Bayonet", "Huntsman Knife", "Falchion Knife",
		"Bowie Knife", "Butterfly Knife", "Shadow Daggers", "Paracord Knife", "Survival Knife", "Nomad Knife",
		"Skeleton Knife", "Stiletto Knife", "Ursus Knife", "Talon Knife", "Classic Knife", "Canis Knife",
		"Outdoor Knife", "Cord Knife", "Gypsy Jackknife", "Widowmaker Knife", "Kukri Knife",
	},
}

// PrintWeaponFireCounts prints how many times each player fired each weapon, sorted by category.
func PrintWeaponFireCounts() {
	fmt.Println("Weapon Fire Counts Per Player (Sorted by Category):")
	fmt.Println("------------------------------------------------------")
	for _, stats := range models.PlayerStatsMap {
		fmt.Printf("Player: %s\n", stats.Name)
		if len(stats.WeaponFireCounts) == 0 {
			fmt.Println("  No weapon fires recorded.")
			continue
		}

		// Track which weapons have already been printed
		printed := make(map[string]bool)

		// Print by category order
		for _, category := range weaponCategories {
			for _, weapon := range category {
				if count, ok := stats.WeaponFireCounts[weapon]; ok {
					fmt.Printf("  %-20s: %d\n", weapon, count)
					printed[weapon] = true
				}
			}
		}

		// Print any remaining weapons not in the categories, sorted alphabetically
		var others []string
		for weapon := range stats.WeaponFireCounts {
			if !printed[weapon] {
				others = append(others, weapon)
			}
		}
		sort.Strings(others)
		for _, weapon := range others {
			fmt.Printf("  %-20s: %d\n", weapon, stats.WeaponFireCounts[weapon])
		}
		fmt.Println()
	}
}
