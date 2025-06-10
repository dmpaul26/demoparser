package main

import (
	"fmt"
	"log"
	"os"

	eventHandlers "demoparser/eventHandlers"
	printers "demoparser/printers"
	utils "demoparser/utils"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func main() {
	var files []string
	foldersToDelete := []string{}

	utils.LoadDemos(&foldersToDelete, &files)

	for _, demoFile := range files {
		parseDemo(demoFile)
	}

	printers.PrintStats()

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

	err = p.ParseToEnd()
	if err != nil {
		log.Panicf("Failed to parse demo %s: %v\n", demoPath, err)
	}
}
