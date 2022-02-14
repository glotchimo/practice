package main

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"
)

func TestAutocompleteHandler(t *testing.T) {
	b := []byte(`{
        "string": "de",
        "dictionary": ["dog", "deer", "deal"]
    }`)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	AutocompleteHandler(rec, req)
	res := rec.Result()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	if bytes.Compare([]byte(`["deer","deal"]`), bytes.Trim(data, "\n")) != 0 {
		t.Errorf("expected\n%b\ngot\n%b", []byte(`["deer","deal"]`), bytes.Trim(data, "\n"))
	}
}
