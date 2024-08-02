package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

func post_test(w http.ResponseWriter, r *http.Request) {
	// Create post end point
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Post request received: %s", body)
}

func post_to_server(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		http.Error(w, "Error marshalling to json:", http.StatusInternalServerError)
		return
	}

	clientName, clientNameExists := jsonData["client"]
	projectName, projectNameExists := jsonData["project"]

	if !clientNameExists || !projectNameExists {
		http.Error(w, "client or project name not found", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Client name: %s\n Project name: %s", clientName, projectName)
}

func main() {
	port := ":8080"

	http.HandleFunc("/get_test", func(w http.ResponseWriter, r *http.Request) {
		out, err := exec.Command("ls", "-l").Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "files: %s", out)
	})

	http.HandleFunc("/post_test", post_test)
	http.HandleFunc("/post", post_to_server)

	log.Printf("Server started %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
