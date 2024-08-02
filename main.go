package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

func create_folder(data map[string]interface{}) {
	return
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
	subsidiary, subsidiaryExists := jsonData["subsidiary"]
	formType, formTypeExists := jsonData["form_type"]

	if !clientNameExists || !projectNameExists || !subsidiaryExists || !formTypeExists {
		var missing string = ""

		if !clientNameExists {
			missing += "client "
		}
		if !projectNameExists {
			missing += "project "
		}
		if !subsidiaryExists {
			missing += "subsidiary "
		}
		if !formTypeExists {
			missing += "form_type "
		}

		http.Error(w, fmt.Sprintf("Missing the following json fields: %s", missing), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Client name:%s,Project name:%s,Subsidiary:%s,Form Type:%s", clientName, projectName, subsidiary, formType)
	// Create folder on the server
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

	http.HandleFunc("/post", post_to_server)

	log.Printf("Server started %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
