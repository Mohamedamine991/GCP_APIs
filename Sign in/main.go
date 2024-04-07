package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
)
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

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password") // Direct plaintext comparison

	// Retrieve the stored password from the database based on the email
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "The email does not exist or the password is incorrect", http.StatusUnauthorized)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	// Directly compare the plaintext passwords
	if password != storedPassword {
		http.Error(w, "The email does not exist or the password is incorrect", http.StatusUnauthorized)
		return
	}

	// If the passwords match, the user is authenticated
	fmt.Fprintf(w, "Sign in successful")
	type UserInfo struct {
		Username         string `json:"username"`
		Email            string `json:"email"`
		GCPProjectName   string `json:"gcp_project_name"`
		AzureClientID    string `json:"azure_client_id"`
		AzureClientSecret string `json:"azure_client_secret"`
		AzureTenantID    string `json:"azure_tenant_id"`
		IBMApiKey        string `json:"ibm_api_key"`
	}
	var user UserInfo
err = db.QueryRow("SELECT username, email, gcp_project_name, azure_client_id, azure_client_secret, azure_tenant_id, ibm_api_key FROM users WHERE email = ?", email).Scan(&user.Username, &user.Email, &user.GCPProjectName, &user.AzureClientID, &user.AzureClientSecret, &user.AzureTenantID, &user.IBMApiKey)
if err != nil {
    http.Error(w, "Server error", http.StatusInternalServerError)
	
    return
}
w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(user)
    if err != nil {
        http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
        return
    }
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/signin", signInHandler)

	fmt.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}