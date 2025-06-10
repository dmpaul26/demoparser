package parser

import (
	"fmt"
	"log"
	"os"

	eventHandlers "demoparser/eventHandlers"

	dem "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

func ParseDemo(demoPath string) {
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
