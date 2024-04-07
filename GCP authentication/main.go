package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func connectGCP(ctx context.Context) *storage.Client {
	fmt.Println("Credentials path:", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	if err != nil {
		log.Fatalf("Failed to create GCP client: %v", err)
	}
	log.Println("Successfully authenticated with GCP")
	return client
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()

	// Authenticate to GCP
	gcpClient := connectGCP(ctx)
	fmt.Print((gcpClient))
	defer gcpClient.Close()

}
