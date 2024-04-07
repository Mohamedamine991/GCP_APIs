package main

import (
	"context"
	"fmt"
	"log"

	compute "cloud.google.com/go/compute/apiv1"
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
	ctx := context.Background()
	projectID := "sound-habitat-418811" // Replace with your GCP project ID
	networkName := "mynetwork"          // Replace with the name of the VPC network you want to delete

	if err := deleteVPCNetwork(ctx, projectID, networkName); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Network %s deleted successfully.\n", networkName)
	}
}
