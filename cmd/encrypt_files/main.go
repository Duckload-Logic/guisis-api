package main

import (
	"context"
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/storage"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Simulate migration without modifying files")
	keyStr := flag.String("key", "", "Storage encryption key (32 bytes hex, base64, or raw). Defaults to STORAGE_ENCRYPTION_KEY env var")
	uploadDir := flag.String("upload-dir", "", "Path to local upload directory. Defaults to UPLOAD_DIR env var or ./uploads")
	storageDriver := flag.String("driver", "", "Storage driver: 'local' or 'lightsail'. Defaults to STORAGE_DRIVER env var")
	bucketName := flag.String("bucket", "", "Lightsail bucket name. Defaults to LIGHTSAIL_BUCKET_NAME env var")
	region := flag.String("region", "", "Lightsail region. Defaults to LIGHTSAIL_REGION env var")

	flag.Parse()

	if *keyStr == "" {
		*keyStr = os.Getenv("STORAGE_ENCRYPTION_KEY")
	}
	if *keyStr == "" {
		*keyStr = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
		log.Println("[INFO] No key provided; using default development key.")
	}

	keyBytes, err := storage.ParseKey(*keyStr)
	if err != nil {
		log.Fatalf("[FATAL] Invalid encryption key: %v", err)
	}

	if *storageDriver == "" {
		*storageDriver = os.Getenv("STORAGE_DRIVER")
		if *storageDriver == "" {
			*storageDriver = "local"
		}
	}

	log.Println("==================================================")
	log.Printf(" Starting File Encryption Migration Tool")
	log.Printf(" Driver:  %s", *storageDriver)
	log.Printf(" Dry Run: %t", *dryRun)
	log.Println("==================================================")

	ctx := context.Background()

	if *storageDriver == "lightsail" {
		if *bucketName == "" {
			*bucketName = os.Getenv("LIGHTSAIL_BUCKET_NAME")
		}
		if *region == "" {
			*region = os.Getenv("LIGHTSAIL_REGION")
		}
		if *bucketName == "" || *region == "" {
			log.Fatalf("[FATAL] Bucket name and region are required for lightsail driver")
		}

		rawStorage, err := storage.NewLightsailStorage(ctx, *bucketName, *region)
		if err != nil {
			log.Fatalf("[FATAL] Failed to initialize Lightsail storage: %v", err)
		}

		log.Printf("[INFO] Migrating files in Lightsail bucket: %s (%s)", *bucketName, *region)
		// Lightsail migration via storage interface or direct keys
		log.Printf("[INFO] Note: for Lightsail S3 storage, run with database records or provide file list.")
		_ = rawStorage
		return
	}

	// Local filesystem migration
	if *uploadDir == "" {
		*uploadDir = os.Getenv("UPLOAD_DIR")
		if *uploadDir == "" {
			*uploadDir = "./uploads"
		}
	}

	cleanDir, err := filepath.Abs(*uploadDir)
	if err != nil {
		log.Fatalf("[FATAL] Invalid upload directory path %q: %v", *uploadDir, err)
	}

	if _, err := os.Stat(cleanDir); os.IsNotExist(err) {
		log.Printf("[WARN] Upload directory %q does not exist. Nothing to encrypt.", cleanDir)
		return
	}

	log.Printf("[INFO] Scanning upload directory: %s", cleanDir)

	var (
		scannedCount   int
		encryptedCount int
		alreadyCount   int
		errorCount     int
	)

	err = filepath.WalkDir(cleanDir, func(filePath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			log.Printf("[ERROR] Error reading %s: %v", filePath, walkErr)
			errorCount++
			return nil
		}

		if d.IsDir() {
			return nil
		}

		scannedCount++
		rel, _ := filepath.Rel(cleanDir, filePath)

		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("[ERROR] Failed to read %s: %v", rel, err)
			errorCount++
			return nil
		}

		if storage.IsEncrypted(data) {
			alreadyCount++
			return nil
		}

		if *dryRun {
			log.Printf("[DRY-RUN] Would encrypt: %s (%d bytes)", rel, len(data))
			encryptedCount++
			return nil
		}

		encryptedData, err := storage.EncryptData(keyBytes, data)
		if err != nil {
			log.Printf("[ERROR] Failed to encrypt %s: %v", rel, err)
			errorCount++
			return nil
		}

		// Write to temporary file first then rename for atomic replacement
		tmpPath := filePath + ".enc_tmp"
		if err := os.WriteFile(tmpPath, encryptedData, 0o600); err != nil {
			log.Printf("[ERROR] Failed to write temporary encrypted file for %s: %v", rel, err)
			errorCount++
			return nil
		}

		if err := os.Rename(tmpPath, filePath); err != nil {
			_ = os.Remove(tmpPath)
			log.Printf("[ERROR] Failed to atomic replace %s: %v", rel, err)
			errorCount++
			return nil
		}

		log.Printf("[ENCRYPTED] %s (%d bytes -> %d bytes)", rel, len(data), len(encryptedData))
		encryptedCount++
		return nil
	})

	if err != nil {
		log.Fatalf("[FATAL] WalkDir failure: %v", err)
	}

	log.Println("==================================================")
	log.Printf(" Migration Summary")
	log.Printf(" Total Scanned:         %d", scannedCount)
	log.Printf(" Already Encrypted:     %d", alreadyCount)
	log.Printf(" Newly Encrypted:       %d", encryptedCount)
	log.Printf(" Errors:                %d", errorCount)
	log.Println("==================================================")
}


