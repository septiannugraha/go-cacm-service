package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/septiannugraha/go-cacm-service/internal/database"
	"github.com/septiannugraha/go-cacm-service/internal/packager"
	"github.com/septiannugraha/go-cacm-service/internal/uploader"
)

// Config represents the situwassa.conf structure
type Config struct {
	Server             string   `json:"server"`
	IntegratedSecurity bool     `json:"integrated_security"`
	UserID             string   `json:"user_id"`
	Password           string   `json:"password"`
	Databases          []string `json:"databases"`
}

const (
	serviceURL = "http://89.116.32.187:4080/upload"
	timeout    = 15 * time.Minute
)

func main() {
	// Load configuration
	config, err := loadConfig("situwassa.conf")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create packager and uploader
	pkg := packager.NewPackager()
	upl := uploader.NewUploader(serviceURL, timeout)

	// Process each database
	for _, dbName := range config.Databases {
		fmt.Printf("Processing database: %s\n", dbName)
		
		// Extract year from database name (e.g., "PEMDA_2024" -> "2024")
		year := extractYear(dbName)
		
		// Connect to database
		client, err := database.NewMSSQLClient(
			config.Server,
			dbName,
			config.UserID,
			config.Password,
			config.IntegratedSecurity,
		)
		if err != nil {
			log.Printf("Failed to connect to %s: %v", dbName, err)
			continue
		}
		defer client.Close()

		// Process different report types
		reportTypes := []struct {
			jenis    string
			name     string
			fileName string
		}{
			{"1", "Revenue", fmt.Sprintf("%s_revenue.bin", dbName)},
			{"2", "Expenses", fmt.Sprintf("%s_expenses.bin", dbName)},
			{"3", "Financing", fmt.Sprintf("%s_financing.bin", dbName)},
		}

		var filesToUpload []string

		for _, report := range reportTypes {
			fmt.Printf("  Fetching %s data...\n", report.name)
			
			// Get data from database
			data, err := client.GetSummaryData(year, report.jenis)
			if err != nil {
				log.Printf("Failed to get %s data: %v", report.name, err)
				continue
			}
			
			fmt.Printf("  Found %d records\n", len(data))
			
			// Package to protobuf
			err = pkg.PackageQueryResults(data, report.fileName)
			if err != nil {
				log.Printf("Failed to package %s data: %v", report.name, err)
				continue
			}
			
			filesToUpload = append(filesToUpload, report.fileName)
		}

		// Upload all files for this database
		fmt.Printf("Uploading %d files...\n", len(filesToUpload))
		err = upl.UploadFiles(filesToUpload)
		if err != nil {
			log.Printf("Failed to upload files: %v", err)
		}

		// Clean up files
		for _, file := range filesToUpload {
			os.Remove(file)
		}
	}

	fmt.Println("Sync completed!")
}

func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func extractYear(dbName string) string {
	// Extract year from database name like "PEMDA_2024"
	if len(dbName) >= 4 {
		possibleYear := dbName[len(dbName)-4:]
		// Simple validation - check if it's a valid year
		for _, ch := range possibleYear {
			if ch < '0' || ch > '9' {
				return "2024" // default
			}
		}
		return possibleYear
	}
	return "2024" // default
}