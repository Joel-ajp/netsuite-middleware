package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

func get_directories(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var directories []string
	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, strings.ToLower(file.Name()))
		}
	}
	return directories, nil
}

func findBestMatch(clientName string, directories []string) string {
	// split by :
	clientName = strings.ToLower(clientName)
	bestMatch := ""
	minDistance := len(clientName) + 1

	for _, dir := range directories {
		distance := levenshtein.DistanceForStrings([]rune(clientName), []rune(dir), levenshtein.DefaultOptions)
		if distance < minDistance {
			minDistance = distance
			bestMatch = dir
		}
	}

	return bestMatch
}

func get_client(client string, root string) (string, error) {
	clients, err := get_directories(root)
	if err != nil {
		return "", err
	}

	bestMatch := findBestMatch(client, clients)
	return bestMatch, nil
}

func create_folder(jsondata map[string]interface{}) (string, error) {
	ROOT_PATH := `/mnt/Survey/Projects`
	TEMPLATE_PATH := `/mnt/Survey/Standards_Templates/8.0 DFS Directory Structure/FE Folder Structure/`
	// add check for midstream, f&e, and upstream

	client, err := get_client(jsondata["client"].(string), ROOT_PATH)
	if err != nil {
		return "", err
	}
	project := jsondata["project"].(string)

	var dest_path string
	if client == "" {
		dest_path = fmt.Sprintf("%s/%s", ROOT_PATH, project)
	} else {
		dest_path = fmt.Sprintf("%s/%s/%s", ROOT_PATH, client, project)
	}
	if err != nil {
		return "", err
	}

	cmd := exec.Command("rsync", "-a", TEMPLATE_PATH, dest_path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error creating folder: %s Output: %s", err, string(output))
	}

	formatted_path := strings.Replace(dest_path, "/mnt/Survey/", "S:/", 1)
	formatted_path = strings.ReplaceAll(formatted_path, "/", "\\")
	return formatted_path, nil
}

func check_json(jsondata map[string]interface{}) string {
	_, clientNameExists := jsondata["client"]
	_, projectNameExists := jsondata["project"]
	_, subsidiaryExists := jsondata["subsidiary"]
	_, formTypeExists := jsondata["form_type"]

	var missingElements []string
	if !clientNameExists {
		missingElements = append(missingElements, "client")
	}
	if !projectNameExists {
		missingElements = append(missingElements, "project")
	}
	if !subsidiaryExists {
		missingElements = append(missingElements, "subsidiary")
	}
	if !formTypeExists {
		missingElements = append(missingElements, "form_type")
	}
	missing := strings.Join(missingElements, " ")

	return missing
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
		http.Error(w, "Error unmarshalling to json:", http.StatusInternalServerError)
		return
	}

	missing := check_json(jsonData)
	if missing != "" {
		http.Error(w, fmt.Sprintf("Missing the following json fields: %s", missing), http.StatusBadRequest)
		return
	} else {
		output, err := create_folder(jsonData)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating folder: %s", err), http.StatusInternalServerError)
			return
		} else {
			fmt.Fprintf(w, "Folder created successfully. Path: %s", output)
		}
	}
}

func main() {
	port := ":8080"
	http.HandleFunc("/post", post_to_server)

	log.Printf("Server started %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
