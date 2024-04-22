package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func deleteSQLInstance(ctx context.Context, projectID, instanceID string) error {
	service, err := sqladmin.NewService(ctx)
	if err != nil {
		return fmt.Errorf("sqladmin.NewService: %v", err)
	}

	if _, err := service.Instances.Delete(projectID, instanceID).Context(ctx).Do(); err != nil {
		return fmt.Errorf("Instances.Delete: %v", err)
	}

	fmt.Printf("SQL instance %s deleted.\n", instanceID)
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	projectID := os.Getenv("PROJECT_ID")
	instanceID := os.Getenv("INSTANCE_ID")
	ctx := context.Background()
	if err := deleteSQLInstance(ctx, projectID, instanceID); err != nil {
		log.Fatal(err)
	}
}
