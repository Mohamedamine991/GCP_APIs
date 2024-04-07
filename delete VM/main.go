package main

import (
	"context"
	"fmt"
	"log"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func deleteComputeInstance(ctx context.Context, projectID, zone, instanceName string) error {
	c, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer c.Close()

	req := &computepb.DeleteInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	op, err := c.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete instance: %v", err)
	}

	fmt.Printf("Delete operation on instance %s: %+v\n", instanceName, op)
	return nil
}

func main() {
	ctx := context.Background()
	if err := deleteComputeInstance(ctx, "sound-habitat-418811", "northamerica-northeast2-a", "amine"); err != nil {
		log.Fatal(err)
	}
}
