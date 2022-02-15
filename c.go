package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func FileHandler(w http.ResponseWriter, r *http.Request) {
	m, h, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Couldn't parse file: %s", err.Error())
		return
	}
	defer m.Close()

	f, err := os.Create(h.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Couldn't create file: %s", err.Error())
	}
	defer f.Close()

	if _, err := io.Copy(f, m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Couldn't copy file contents: %s", err.Error())
		return
	}
}

func c() {
	http.HandleFunc("/", FileHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
