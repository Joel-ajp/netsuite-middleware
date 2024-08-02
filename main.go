package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		out, err := exec.Command("ls", "-l").Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "files: %s", out)
	})

	log.Print("Server started")

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
