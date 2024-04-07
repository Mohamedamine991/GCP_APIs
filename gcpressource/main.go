package main

import (
	"context"
	"fmt"
	"log"

	storage "cloud.google.com/go/storage"
	compute "cloud.google.com/go/compute/apiv1"
	"google.golang.org/api/iterator"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
		"google.golang.org/api/cloudresourcemanager/v1"

)
func listProjectIamPolicies(ctx context.Context, projectID string) {
	// Create a service object for the Cloud Resource Manager API
	service, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		log.Fatalf("cloudresourcemanager.NewService: %v", err)
	}

	// Use the getIamPolicy method to retrieve IAM policies for the project
	resp, err := service.Projects.GetIamPolicy(projectID, &cloudresourcemanager.GetIamPolicyRequest{}).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Projects.GetIamPolicy: %v", err)
	}

	fmt.Printf("IAM Policy for project %s:\n", projectID)
	for _, binding := range resp.Bindings {
		fmt.Printf("Role: %s\n", binding.Role)
		for _, member := range binding.Members {
			fmt.Printf(" - Member: %s\n", member)
		}
	}
}
func listInstances(ctx context.Context, projectID string) {
	c, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create instance client: %v", err)
	}
	defer c.Close()

	req := &computepb.AggregatedListInstancesRequest{
		Project: projectID,
	}

	it := c.AggregatedList(ctx, req)
	fmt.Println("Compute Engine Instances:")
	for {
		pair, err := it.Next()
		if err != nil {
			break // End of list
		}
		if pair.Value.Instances != nil {
			for _, instance := range pair.Value.Instances {
				fmt.Printf("- %s (Zone: %s, Machine Type: %s)\n", instance.GetName(), instance.GetZone(), instance.GetMachineType())
			}
		}
	}
	if err != nil {
		log.Fatalf("Failed to list instances: %v", err)
	}
}

func listStorageBuckets(ctx context.Context, projectID string) error {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer storageClient.Close()

	it := storageClient.Buckets(ctx, projectID)
	fmt.Println("Cloud Storage Buckets:")
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("- %s (Location: %s)\n", attrs.Name, attrs.Location)
	}
	return nil
}

func listVPCNetworks(ctx context.Context, projectID string) error {
	c, err := compute.NewNetworksRESTClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Networks client: %v", err)
	}
	defer c.Close()

	req := &computepb.ListNetworksRequest{
		Project: projectID,
	}

	it := c.List(ctx, req)
	fmt.Println("VPC Networks:")
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Failed to list networks: %v", err)
		}
		fmt.Printf("- %s (Auto Create Subnetworks: %t)\n", resp.GetName(), resp.GetAutoCreateSubnetworks())
	}
	return nil
}

func main() {
	ctx := context.Background()
	projectID := "sound-habitat-418811" // Replace with your actual project ID

	listInstances(ctx, projectID)

	if err := listStorageBuckets(ctx, projectID); err != nil {
		log.Printf("Error listing Cloud Storage buckets: %v", err)
	}

	if err := listVPCNetworks(ctx, projectID); err != nil {
		log.Printf("Error listing VPC Networks: %v", err)
	}

		listProjectIamPolicies(ctx, projectID)

}
