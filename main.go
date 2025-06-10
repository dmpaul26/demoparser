package main

import (
	"fmt"
	"log"
	"os"

	parser "demoparser/parser"
	printers "demoparser/printers"
	utils "demoparser/utils"
)

func main() {
	var files []string
	foldersToDelete := []string{}

	utils.LoadDemos(&foldersToDelete, &files)

	for _, demoFile := range files {
		parser.ParseDemo(demoFile)
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
