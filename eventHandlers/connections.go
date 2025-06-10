package eventHandlers

import (
	"demoparser/models"
	"fmt"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

var DebugConnections = false

// HandlePlayerConnectEvent logs the name of the player who connects.
func HandlePlayerConnectEvent(parser dem.Parser, event events.PlayerConnect) {
	if event.Player != nil {
		if DebugConnections {
			fmt.Printf("Player connected: %s at tick %d\n", event.Player.Name, parser.GameState().IngameTick())
		}

		models.TryInitializePlayerStatsMap(event.Player.SteamID64, event.Player.Name)
	}
}

// HandlePlayerConnectEvent logs the name of the player who connects.
func HandlePlayerDisconnectedEvent(parser dem.Parser, event events.PlayerDisconnected) {
	if event.Player != nil {
		if DebugConnections {
			fmt.Printf("Player disconnected: %s at tick %d\n", event.Player.Name, parser.GameState().IngameTick())
		}

		models.TryInitializePlayerStatsMap(event.Player.SteamID64, event.Player.Name)
	}
}
