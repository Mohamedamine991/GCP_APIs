package main

import (
	"context"
	"fmt"
	"log"
	"os"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/joho/godotenv"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

// deleteVPCNetwork deletes the specified VPC network in the given project.
func deleteVPCNetwork(ctx context.Context, projectID, networkName string) error {
	c, err := compute.NewNetworksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewNetworksRESTClient: %v", err)
	}
	defer c.Close()

	req := &computepb.DeleteNetworkRequest{
		Project: projectID,
		Network: networkName,
	}

	op, err := c.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete network: %v", err)
	}

	fmt.Printf("Delete operation on network %s: %+v\n", networkName, op)
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	networkName := os.Getenv("NETWORK_NAME") // Replace with the name of the VPC network you want to delete

	if err := deleteVPCNetwork(ctx, projectID, networkName); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Network %s deleted successfully.\n", networkName)
	}
}
