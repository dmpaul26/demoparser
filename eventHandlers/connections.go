package eventHandlers

import (
	"demoparser/models"
	"fmt"

	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

var DebugConnections = false

// HandlePlayerConnectEvent logs the name of the player who connects.
func HandlePlayerConnectEvent(event events.PlayerConnect) {
	if event.Player != nil {
		if DebugConnections {
			fmt.Printf("Player connected: %s\n", event.Player.Name)
		}

		models.TryInitializePlayerStatsMap(event.Player.SteamID64, event.Player.Name)
	}
}

// HandlePlayerConnectEvent logs the name of the player who connects.
func HandlePlayerDisconnectedEvent(event events.PlayerDisconnected) {
	if event.Player != nil {
		if DebugConnections {
			fmt.Printf("Player disconnected: %s\n", event.Player.Name)
		}

		models.TryInitializePlayerStatsMap(event.Player.SteamID64, event.Player.Name)
	}
}
