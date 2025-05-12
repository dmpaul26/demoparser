package utils

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	zstd "github.com/klauspost/compress/zstd"
)

// extractGzFile extracts a .gz file to the specified output directory.
func ExtractGzFile(gzFilePath string, outputDir string) error {
	// Open the .gz file
	gzFile, err := os.Open(gzFilePath)
	if err != nil {
		return fmt.Errorf("failed to open .gz file: %v", err)
	}
	defer gzFile.Close()

	// Create a gzip reader
	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	// Extract the file name from the .gz file
	extractedFileName := strings.TrimSuffix(filepath.Base(gzFilePath), ".gz")
	extractedFilePath := filepath.Join(outputDir, extractedFileName)

	// Create the output file
	outputFile, err := os.Create(extractedFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Copy the contents of the .gz file to the output file
	_, err = io.Copy(outputFile, gzReader)
	if err != nil {
		return fmt.Errorf("failed to extract .gz file: %v", err)
	}

	return nil
}

// ExtractZstFile extracts a .zst file to the specified output directory.
func ExtractZstFile(zstFilePath string, outputDir string) error {
	// Open the .zst file
	zstFile, err := os.Open(zstFilePath)
	if err != nil {
		return fmt.Errorf("failed to open .zst file: %v", err)
	}
	defer zstFile.Close()

	// Create a zstd decoder
	zstDecoder, err := zstd.NewReader(zstFile)
	if err != nil {
		return fmt.Errorf("failed to create zstd decoder: %v", err)
	}
	defer zstDecoder.Close()

	// Extract the file name from the .zst file
	extractedFileName := strings.TrimSuffix(filepath.Base(zstFilePath), ".zst")
	extractedFilePath := filepath.Join(outputDir, extractedFileName)

	// Create the output file
	outputFile, err := os.Create(extractedFilePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Copy the contents of the .zst file to the output file
	_, err = io.Copy(outputFile, zstDecoder)
	if err != nil {
		return fmt.Errorf("failed to extract .zst file: %v", err)
	}

	return nil
}
