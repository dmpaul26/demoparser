package eventHandlers

import (
	"fmt"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

var DebugChat = false

// HandleChatMessage writes the player's name in brackets, a colon and space, then the message.
// If the message is team-only, "TEAM" is added at the beginning.
func HandleChatMessage(parser dem.Parser, event events.ChatMessage) {
	if DebugChat {
		if event.IsChatAll {
			fmt.Printf("TEAM [%s]: %s\n", event.Sender.Name, event.Text)
		} else {
			fmt.Printf("[%s]: %s\n", event.Sender.Name, event.Text)
		}
	}
}
