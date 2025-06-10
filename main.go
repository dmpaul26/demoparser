package main

import (
	"fmt"
	"log"
	"os"

	"demoparser/eventHandlers"
	"demoparser/printers"
	"demoparser/utils"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func debugToggles() {
	// Uncomment the following lines to enable debug output for specific events
	//eventHandlers.DebugChat = true
	//eventHandlers.DebugRounds = true
	//eventHandlers.DebugRoundEnds = true
	//eventHandlers.DebugConnections = true
	//eventHandlers.DebugWeaponFire = true
	//eventHandlers.WeaponFirePlayerFilter = "" // Set to "" to log all, or a name to filter
	//eventHandlers.DebugSpotters = true
	//models.DebugPlayerStatsInit = true // or false to disable
	//eventHandlers.MaxRoundDebug = 1
}

func printerToggles() {
	// Uncomment the following lines to enable specific printers
	printers.PrintStats()
	//printers.PrintWeaponFireCounts()
}

func main() {
	var files []string
	foldersToDelete := []string{}

	debugToggles()
	utils.LoadDemos(&foldersToDelete, &files)

	//files = files[:1] // *** Limit to 1 demo for testing
	for _, demoFile := range files {
		parseDemo(demoFile)
		resetGlobals()
	}

	printerToggles()

	cleanupFolders(foldersToDelete)
	fmt.Println("Demo processing completed.")
}

func cleanupFolders(foldersToDelete []string) {
	for _, folder := range foldersToDelete {
		err := os.RemoveAll(folder)
		if err != nil {
			log.Printf("Failed to delete folder %s: %v\n", folder, err)
		} else {
			fmt.Printf("Deleted folder: %s\n", folder)
		}
	}
}

func parseDemo(demoPath string) {
	f, err := os.Open(demoPath)
	if err != nil {
		log.Panicf("Failed to open demo file %s: %v\n", demoPath, err)
	}
	defer f.Close()

	parser := demoinfocs.NewParser(f)
	defer parser.Close()

	fmt.Printf("Processing demo: %s\n", demoPath)

	parser.RegisterEventHandler(func(event events.Kill) { eventHandlers.HandleKillEvent(parser, event) })
	parser.RegisterEventHandler(func(event events.PlayerHurt) { eventHandlers.HandlePlayerHurtEvent(parser, event) })
	parser.RegisterEventHandler(func(event events.WeaponFire) { eventHandlers.HandleWeaponFireEvent(parser, event) })
	parser.RegisterEventHandler(func(event events.PlayerSpottersChanged) {
		eventHandlers.HandlePlayerSpottersChangedEvent(parser, event)
	})
	parser.RegisterEventHandler(func(event events.RoundStart) { eventHandlers.HandleRoundStartEvent(parser, event) })
	parser.RegisterEventHandler(func(event events.RoundEnd) { eventHandlers.HandleRoundEndEvent(parser, event) })
	parser.RegisterEventHandler(func(event events.MatchStart) { eventHandlers.HandleMatchStartedEvent(parser, event) })
	parser.RegisterEventHandler(func(event events.AnnouncementFinalRound) { eventHandlers.HandleFinalRoundEvent(parser, event) })
	parser.RegisterEventHandler(func(event events.AnnouncementWinPanelMatch) { eventHandlers.HandleWinPanelMatchEvent(parser, event) })
	parser.RegisterEventHandler(func(event events.PlayerConnect) { eventHandlers.HandlePlayerConnectEvent(parser, event) })
	parser.RegisterEventHandler(func(event events.PlayerDisconnected) { eventHandlers.HandlePlayerDisconnectedEvent(parser, event) })
	parser.RegisterEventHandler(func(event events.ChatMessage) { eventHandlers.HandleChatMessage(parser, event) })

	err = parser.ParseToEnd()
	if err != nil {
		log.Panicf("Failed to parse demo %s: %v\n", demoPath, err)
	}
}

func resetGlobals() {
	// Reset global variables or states if needed
	eventHandlers.RoundCount = 0
	eventHandlers.MatchStartCount = 0
}
