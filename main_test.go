package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubstringLengthHandler(t *testing.T) {
	b := bytes.Buffer{}
	b.WriteString("3,abcabcdefdefdef")

	req := httptest.NewRequest(http.MethodPost, "/?=2", &b)
	rec := httptest.NewRecorder()

	SubstringLengthHandler(rec, req)
	res := rec.Result()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Errorf("expected 200 got %d", res.StatusCode)
	}

	if string(data) != "9" {
		t.Errorf("expected 9 got %s", string(data))
	}
}
