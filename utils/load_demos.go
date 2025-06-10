package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// handleArchive extracts an archive (.gz or .zst), finds .dem files, and updates foldersToDelete and files.
func handleArchive(inputFile string, extractFunc func(string, string) error, foldersToDelete *[]string, files *[]string) {
	outputDir := strings.SplitN(filepath.Base(inputFile), ".", 2)[0]
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Panicf("Failed to create output directory: %v\n", err)
	}

	err = extractFunc(inputFile, outputDir)
	if err != nil {
		log.Panicf("Failed to extract archive: %v\n", err)
	}

	extractedFiles, err := filepath.Glob(filepath.Join(outputDir, "*.dem"))
	*foldersToDelete = append(*foldersToDelete, outputDir)
	if err != nil {
		log.Panicf("Failed to read extracted files: %v\n", err)
	}

	if len(extractedFiles) > 0 {
		*files = append(*files, extractedFiles...)
	} else {
		fmt.Println("No .dem files found in the extracted folder.")
	}
}

// LoadDemos processes input files and populates foldersToDelete and files slices.
func LoadDemos(foldersToDelete *[]string, files *[]string) {
	// Check if a file parameter is passed
	if len(os.Args) > 1 {
		for _, inputFile := range os.Args[1:] {
			if strings.HasSuffix(inputFile, ".gz") {
				handleArchive(inputFile, ExtractGzFile, foldersToDelete, files)
			} else if strings.HasSuffix(inputFile, ".zst") {
				handleArchive(inputFile, ExtractZstFile, foldersToDelete, files)
			} else if strings.HasSuffix(inputFile, ".dem") {
				*files = []string{inputFile}
			}
		}
	} else {
		// Default to processing all .dem files in the current directory
		foundFiles, err := filepath.Glob("*.dem")
		if err != nil {
			log.Panic("Failed to read demo files: ", err)
		}

		if len(foundFiles) == 0 {
			fmt.Println("No .dem files found in the current folder.")
			return
		}
		*files = foundFiles
	}
}
