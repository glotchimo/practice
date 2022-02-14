// Implement an autocomplete system. That is, given a query string `s` and a set of all
// possible query strings, return all strings in the set that have s as a prefix.
// For example, given the query string `de` and the set of strings `['dog', 'deer', 'deal']`,
// return `['deer', 'deal']`. Implement as an HTTP endpoint using bluele/gcache to cache completions.

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
)

type AutocompleteRequest struct {
	String     string   `json:"string"`
	Dictionary []string `json:"dictionary"`
}

func autocompleteLinear(s string, d []string) []string {
	matches := []string{}
	for _, w := range d {
		if w[:len(s)] == s {
			matches = append(matches, w)
		}
	}

	return matches
}

func AutocompleteHandler(w http.ResponseWriter, r *http.Request) {
	var body AutocompleteRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Couldn't decode request body: %s", err.Error())
		return
	}
	defer r.Body.Close()

	hasher := sha256.New()
	hasher.Write([]byte(body.String))
	hash := hasher.Sum(nil)
	strHash := string(hash)

	options, errGet := Cache.Get(strHash)
	if errGet != nil {
		options = autocompleteLinear(body.String, body.Dictionary)
		if errSet := Cache.Set(strHash, options); errSet != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Cache write failed: %s", errSet.Error())
			return
		}
	}

	if err := json.NewEncoder(w).Encode(options); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error: %s", err.Error())
		return
	}
}

func B() {
	http.HandleFunc("/", AutocompleteHandler)
	http.ListenAndServe(":8080", nil)
}
