package main

import (
	"context"
	"fmt"
	"log"

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
	ctx := context.Background()
	if err := deleteSQLInstance(ctx, "sound-habitat-418811", "mysql"); err != nil {
		log.Fatal(err)
	}
}
