package eventHandlers

import (
	"fmt"

	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

var DebugChat = false

// HandleChatMessage writes the player's name in brackets, a colon and space, then the message.
// If the message is team-only, "TEAM" is added at the beginning.
func HandleChatMessage(e events.ChatMessage) {
	if DebugChat {
		if e.IsChatAll {
			fmt.Printf("TEAM [%s]: %s\n", e.Sender.Name, e.Text)
		} else {
			fmt.Printf("[%s]: %s\n", e.Sender.Name, e.Text)
		}
	}
}
