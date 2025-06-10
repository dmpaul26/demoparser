package main

import (
	"fmt"
	"log"
	"os"

	eventHandlers "demoparser/eventHandlers"
	utils "demoparser/utils"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func main() {
	var files []string
	foldersToDelete := []string{}

	// toggle-able debug flags
	//eventHandlers.DebugChat = true
	//eventHandlers.DebugRounds = true
	//eventHandlers.DebugRoundEnds = true
	//eventHandlers.DebugConnections = true

	utils.LoadDemos(&foldersToDelete, &files)

	for _, demoFile := range files {
		parseDemo(demoFile)
		resetGlobals()
	}

	//printers.PrintStats()

	// Clean up extracted folders
	for _, folder := range foldersToDelete {
		err := os.RemoveAll(folder)
		if err != nil {
			log.Printf("Failed to delete folder %s: %v\n", folder, err)
		} else {
			fmt.Printf("Deleted folder: %s\n", folder)
		}
	}
	fmt.Println("Demo processing completed.")
}

func parseDemo(demoPath string) {
	f, err := os.Open(demoPath)
	if err != nil {
		log.Panicf("Failed to open demo file %s: %v\n", demoPath, err)
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()

	fmt.Printf("Processing demo: %s\n", demoPath)

	p.RegisterEventHandler(func(e events.Kill) { eventHandlers.HandleKillEvent(p, e) })
	p.RegisterEventHandler(func(e events.PlayerHurt) { eventHandlers.HandlePlayerHurtEvent(p, e) })
	p.RegisterEventHandler(func(e events.WeaponFire) { eventHandlers.HandleWeaponFireEvent(p, e) })
	p.RegisterEventHandler(func(e events.PlayerSpottersChanged) { eventHandlers.HandlePlayerSpottersChangedEvent(p, e) })
	p.RegisterEventHandler(func(e events.RoundStart) { eventHandlers.HandleRoundStartEvent(e) })
	p.RegisterEventHandler(func(e events.RoundEnd) { eventHandlers.HandleRoundEndEvent(e) })
	p.RegisterEventHandler(func(e events.MatchStart) { eventHandlers.HandleMatchStartedEvent(e) })
	p.RegisterEventHandler(func(e events.AnnouncementFinalRound) { eventHandlers.HandleFinalRoundEvent(e) })
	p.RegisterEventHandler(func(e events.AnnouncementWinPanelMatch) { eventHandlers.HandleWinPanelMatchEvent(e) })
	p.RegisterEventHandler(func(e events.PlayerConnect) { eventHandlers.HandlePlayerConnectEvent(e) })
	p.RegisterEventHandler(func(e events.PlayerDisconnected) { eventHandlers.HandlePlayerDisconnectedEvent(e) })
	p.RegisterEventHandler(func(e events.ChatMessage) { eventHandlers.HandleChatMessage(e) })

	err = p.ParseToEnd()
	if err != nil {
		log.Panicf("Failed to parse demo %s: %v\n", demoPath, err)
	}
}

func resetGlobals() {
	// Reset global variables or states if needed
	eventHandlers.RoundCount = 0
	eventHandlers.MatchStartCount = 0
}
