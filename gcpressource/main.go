package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iterator"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

// Handler for listing IAM policies
func listProjectIamPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID") // Assumes PROJECT_ID is set in environment

	service, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating cloudresourcemanager service: %v", err), http.StatusInternalServerError)
		return
	}

	resp, err := service.Projects.GetIamPolicy(projectID, &cloudresourcemanager.GetIamPolicyRequest{}).Context(ctx).Do()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting IAM policy: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// Handler for listing Compute Engine instances
func listInstancesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID") // Assumes PROJECT_ID is set in environment

	c, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating instance client: %v", err), http.StatusInternalServerError)
		return
	}
	defer c.Close()

	req := &computepb.AggregatedListInstancesRequest{Project: projectID}
	it := c.AggregatedList(ctx, req)
	var instances []*computepb.Instance

	for {
		pair, err := it.Next()
		if err != nil {
			break // End of list
		}
		if pair.Value.Instances != nil {
			instances = append(instances, pair.Value.Instances...)
		}
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing instances: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(instances)
}

// Handler for listing Storage Buckets
func listStorageBucketsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID") // Assumes PROJECT_ID is set in environment

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating storage client: %v", err), http.StatusInternalServerError)
		return
	}
	defer storageClient.Close()

	it := storageClient.Buckets(ctx, projectID)
	var buckets []*storage.BucketAttrs

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("Error listing buckets: %v", err), http.StatusInternalServerError)
			return
		}
		buckets = append(buckets, attrs)
	}

	json.NewEncoder(w).Encode(buckets)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	http.HandleFunc("/iam-policies", listProjectIamPoliciesHandler)
	http.HandleFunc("/instances", listInstancesHandler)
	http.HandleFunc("/storage-buckets", listStorageBucketsHandler)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
