package utils

import (
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
)

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
