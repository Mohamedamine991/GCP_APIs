package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
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
	ctx := context.Background()
	if err := deleteStorageBucket(ctx, "stal02"); err != nil {
		log.Fatal(err)
	}
}
