package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// R2Configuration holds the necessary details for R2 connection
type R2Configuration struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	FolderName      string
}

// NeonConfiguration holds the necessary details for Neon DB connection
type NeonConfiguration struct {
	ConnectionString string
	DatabaseName     string
	SchemaName       string
	TableName        string
}

func main() {
	// Load .env file. This should be one of the first things in main.
	err := godotenv.Load() // By default, it loads .env from the current directory
	if err != nil {
		log.Println("No .env file found, relying on environment variables set elsewhere.")
	}

	// --- Configuration ---
	// IMPORTANT: Replace with your actual Cloudflare R2 credentials and details
	r2Config := R2Configuration{
		AccountID:       "b83d3b9a7092a7da1bee02c72968b81e",
		AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),    // Best practice: use environment variables
		SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"), // Best practice: use environment variables
		BucketName:      "dev-highres",
		FolderName:      "004855a7-9d70-5446-bc6d-b7edac9a2479/",
	}

	// IMPORTANT: Replace with your actual Neon connection string and details
	// Example: "postgres://user:password@host:port/dbname?sslmode=require"
	neonConnStr := os.Getenv("NEON_CONNECTION_STRING") // Best practice: use environment variables
	if neonConnStr == "" {
		log.Println("NEON_CONNECTION_STRING environment variable is not set. Using placeholder.")
		neonConnStr = "postgres://neon_user:neon_password@ep-plain-block-00000000.us-east-2.aws.neon.tech/cozy-dev?sslmode=require" // Replace with a valid default or ensure env var is set
	}

	neonConfig := NeonConfiguration{
		ConnectionString: neonConnStr,
		DatabaseName:     "cozy-dev",
		SchemaName:       "public",
		TableName:        "pipeline_defs",
	}

	// --- Cloudflare R2 Interaction ---
	fmt.Println("--- Cloudflare R2 ---")
	startTimeR2 := time.Now()
	files, numFiles, err := listR2Files(r2Config)
	durationR2 := time.Since(startTimeR2)

	if err != nil {
		log.Printf("Error listing files from R2: %v", err)
	} else {
		fmt.Printf("Found %d files in R2 bucket '%s' folder '%s':\n", numFiles, r2Config.BucketName, r2Config.FolderName)
		for _, file := range files {
			fmt.Printf("- %s\n", file)
		}
		fmt.Printf("Time taken to list R2 files: %s\n", durationR2)
	}
	fmt.Println("---------------------")
	fmt.Println()

	// --- Cloudflare R2 Count Only Interaction ---
	fmt.Println("--- Cloudflare R2 (Count Only) ---")
	startTimeR2Count := time.Now()
	numFilesOnly, errCount := countR2FilesOnly(r2Config)
	durationR2Count := time.Since(startTimeR2Count)

	if errCount != nil {
		log.Printf("Error counting files in R2: %v", errCount)
	} else {
		fmt.Printf("Counted %d files in R2 bucket '%s' folder '%s'.\n", numFilesOnly, r2Config.BucketName, r2Config.FolderName)
		fmt.Printf("Time taken to count R2 files: %s\n", durationR2Count)
	}
	fmt.Println("----------------------------------")
	fmt.Println()

	// --- Neon DB Interaction ---
	fmt.Println("--- Neon PostgreSQL ---")
	startTimeNeon := time.Now()
	record, err := queryNeonPostgres(neonConfig)
	durationNeon := time.Since(startTimeNeon)

	if err != nil {
		log.Printf("Error querying Neon PostgreSQL: %v", err)
	} else {
		fmt.Printf("First record from Neon DB '%s'.'%s'.'%s':\n", neonConfig.DatabaseName, neonConfig.SchemaName, neonConfig.TableName)
		// Assuming the record is a map or a struct; adjust printing as needed
		for col, val := range record {
			fmt.Printf("- %s: %v\n", col, val)
		}
		fmt.Printf("Time taken to query Neon PostgreSQL: %s\n", durationNeon)
	}
	fmt.Println("-----------------------")
}

// listR2Files connects to Cloudflare R2 and lists files in a specified bucket and folder.
func listR2Files(cfg R2Configuration) ([]string, int, error) {
	if cfg.AccountID == "YOUR_R2_ACCOUNT_ID" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return nil, 0, fmt.Errorf("R2 credentials are not configured. Please set AccountID, R2_ACCESS_KEY_ID, and R2_SECRET_ACCESS_KEY")
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID),
			HostnameImmutable: true,
			Source:            aws.EndpointSourceCustom,
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
		config.WithRegion("auto"), // R2 specific region
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load R2 configuration: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)

	var files []string
	var continuationToken *string
	totalFiles := 0

	for {
		params := &s3.ListObjectsV2Input{
			Bucket:            aws.String(cfg.BucketName),
			Prefix:            aws.String(cfg.FolderName),
			ContinuationToken: continuationToken,
		}

		resp, err := s3Client.ListObjectsV2(context.TODO(), params)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to list objects in R2 bucket '%s' with prefix '%s': %w", cfg.BucketName, cfg.FolderName, err)
		}

		for _, item := range resp.Contents {
			// Exclude the folder itself if it appears as an object
			if item.Key != nil && *item.Key != cfg.FolderName {
				files = append(files, *item.Key)
				totalFiles++
			}
		}

		if !aws.ToBool(resp.IsTruncated) {
			break // No more objects to list
		}
		continuationToken = resp.NextContinuationToken
	}

	return files, totalFiles, nil
}

// countR2FilesOnly connects to Cloudflare R2 and counts files in a specified bucket and folder.
func countR2FilesOnly(cfg R2Configuration) (int, error) {
	if cfg.AccountID == "YOUR_R2_ACCOUNT_ID" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return 0, fmt.Errorf("R2 credentials are not configured. Please set AccountID, R2_ACCESS_KEY_ID, and R2_SECRET_ACCESS_KEY")
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID),
			HostnameImmutable: true,
			Source:            aws.EndpointSourceCustom,
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
		config.WithRegion("auto"), // R2 specific region
	)
	if err != nil {
		return 0, fmt.Errorf("failed to load R2 configuration: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)

	var continuationToken *string
	totalFiles := 0

	for {
		params := &s3.ListObjectsV2Input{
			Bucket:            aws.String(cfg.BucketName),
			Prefix:            aws.String(cfg.FolderName),
			ContinuationToken: continuationToken,
		}

		resp, err := s3Client.ListObjectsV2(context.TODO(), params)
		if err != nil {
			return 0, fmt.Errorf("failed to list objects in R2 bucket '%s' with prefix '%s': %w", cfg.BucketName, cfg.FolderName, err)
		}

		for _, item := range resp.Contents {
			// Exclude the folder itself if it appears as an object
			if item.Key != nil && *item.Key != cfg.FolderName {
				totalFiles++
			}
		}

		if !aws.ToBool(resp.IsTruncated) {
			break // No more objects to list
		}
		continuationToken = resp.NextContinuationToken
	}

	return totalFiles, nil
}

// queryNeonPostgres connects to a Neon PostgreSQL database and fetches the first record from a specified table.
func queryNeonPostgres(cfg NeonConfiguration) (map[string]interface{}, error) {
	if cfg.ConnectionString == "" || cfg.ConnectionString == "postgres://neon_user:neon_password@ep-plain-block-00000000.us-east-2.aws.neon.tech/cozy-dev?sslmode=require" {
		return nil, fmt.Errorf("Neon connection string is not configured. Please set NEON_CONNECTION_STRING environment variable or update the default")
	}

	db, err := sql.Open("postgres", cfg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open Neon PostgreSQL connection: %w", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping Neon PostgreSQL database: %w", err)
	}

	query := fmt.Sprintf("SELECT * FROM %s.%s LIMIT 1;", pq.QuoteIdentifier(cfg.SchemaName), pq.QuoteIdentifier(cfg.TableName))
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query on Neon PostgreSQL table '%s.%s': %w", cfg.SchemaName, cfg.TableName, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no records found in Neon PostgreSQL table '%s.%s'", cfg.SchemaName, cfg.TableName)
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns from Neon PostgreSQL result: %w", err)
	}

	// Create a slice of interface{}'s to represent each column
	// and a second slice to contain pointers to each item in the first slice
	values := make([]interface{}, len(cols))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	err = rows.Scan(scanArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row from Neon PostgreSQL result: %w", err)
	}

	record := make(map[string]interface{})
	for i, colName := range cols {
		val := values[i]
		// Handle common data types, convert []byte to string for text-based columns
		if b, ok := val.([]byte); ok {
			record[colName] = string(b)
		} else {
			record[colName] = val
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows from Neon PostgreSQL result: %w", err)
	}

	return record, nil
}
