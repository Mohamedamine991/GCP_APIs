package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username         string
	Email            string
	Password         string // Storing plaintext password for the sake of example
	GCPJSON          []byte
	GCPProjectName   string
	AzureClientID    string
	AzureClientSecret string
	AzureTenantID    string
	IBMApiKey        string
}

var db *sql.DB

func initDB() {
	var err error
	// Replace with your actual database connection details
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/platform")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB file size limit
	if err != nil {
		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("gcp_json_file")
	if err != nil {
		http.Error(w, "Error retrieving the GCP JSON file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	var gcpJSON bytes.Buffer
	_, err = io.Copy(&gcpJSON, file)
	if err != nil {
		http.Error(w, "Error reading the GCP JSON file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user := User{
		Username:         r.FormValue("username"),
		Email:            r.FormValue("email"),
		Password:         r.FormValue("password"), // Should be hashed and secured in production
		GCPJSON:          gcpJSON.Bytes(),
		GCPProjectName:   r.FormValue("gcp_project_id"),
		AzureClientID:    r.FormValue("azure_client_id"),
		AzureClientSecret: r.FormValue("azure_client_secret"),
		AzureTenantID:    r.FormValue("azure_tenant_id"),
		IBMApiKey:        r.FormValue("ibm_api_key"),
	}

	// Using a prepared statement for security against SQL injections
	stmt, err := db.Prepare("INSERT INTO users(username, email, password, gcp_json, gcp_project_name, azure_client_id, azure_client_secret, azure_tenant_id, ibm_api_key) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Error preparing SQL statement: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Email, user.Password, user.GCPJSON, user.GCPProjectName, user.AzureClientID, user.AzureClientSecret, user.AzureTenantID, user.IBMApiKey)
	if err != nil {
		http.Error(w, "Error executing SQL statement: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User created successfully")
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/signup", signupHandler)

	fmt.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
