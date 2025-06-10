package eventHandlers

import (
	"fmt"

	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

var DebugConnections = false

// HandlePlayerConnectEvent logs the name of the player who connects.
func HandlePlayerConnectEvent(e events.PlayerConnect) {
	if e.Player != nil {
		if DebugConnections {
			fmt.Printf("Player connected: %s\n", e.Player.Name)
		}
	}
}

// HandlePlayerConnectEvent logs the name of the player who connects.
func HandlePlayerDisconnectedEvent(e events.PlayerDisconnected) {
	if e.Player != nil {
		if DebugConnections {
			fmt.Printf("Player disconnected: %s\n", e.Player.Name)
		}
	}
}
