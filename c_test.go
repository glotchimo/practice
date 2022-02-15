package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFileHandler(t *testing.T) {
	f, err := os.Create("file")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	b := bytes.Buffer{}
	bw := multipart.NewWriter(&b)
	fw, err := bw.CreateFormFile("file", "file")
	if err != nil {
		t.Error(err)
	}
	if _, err := io.Copy(fw, f); err != nil {
		t.Error(err)
	}
	bw.Close()

	req := httptest.NewRequest("POST", "http://localhost:8080", &b)
	req.Header.Set("Content-Type", bw.FormDataContentType())
	rec := httptest.NewRecorder()

	FileHandler(rec, req)
	res := rec.Result()

	if res.StatusCode != 200 {
		t.Error(err)
	}
}
