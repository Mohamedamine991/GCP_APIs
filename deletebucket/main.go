package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
)

func deleteStorageBucket(ctx context.Context, bucketName string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	if err := client.Bucket(bucketName).Delete(ctx); err != nil {
		return fmt.Errorf("Bucket(%q).Delete: %v", bucketName, err)
	}

	fmt.Printf("Bucket %s deleted.\n", bucketName)
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	ctx := context.Background()
	bucketID := os.Getenv("BUCKET_ID")
	if err := deleteStorageBucket(ctx, bucketID); err != nil {
		log.Fatal(err)
	}
}
