package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sgryczan/scanley/models"
	"github.com/sgryczan/scanley/scanner"
)

func TestGetScanHandler(t *testing.T) {
	// set up a temp inventory dir
	targetDir := "scans/"
	id := "test123"
	err := os.Mkdir(targetDir, 0777)
	defer os.RemoveAll(targetDir)

	res := &models.ScanResults{
		Date: "now",
	}
	writer := scanner.NewScanWriter()
	err = writer.Update(res, id)
	if err != nil {
		t.Fatalf("%v", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("/scan/%s", id), nil)

	// inject for testing
	// https://www.gorillatoolkit.org/pkg/mux#SetURLVars
	vars := map[string]string{
		"id": id,
	}
	req = mux.SetURLVars(req, vars)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetScanHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned status code %v want %v", status, http.StatusOK)
	}

	expected, err := json.MarshalIndent(*res, "", "  ")
	if rr.Body.String() != string(expected) {
		t.Errorf("handler returned unexpected body: got %v\n want \n%v",
			rr.Body.String(), expected)
	}
}

func TestGetScanHandlerBad(t *testing.T) {
	// set up a temp inventory dir
	targetDir := "scans/"
	err := os.Mkdir(targetDir, 0777)
	defer os.RemoveAll(targetDir)

	if err != nil {
		t.Fatalf("%v", err)
	}

	id := "idontexist"
	req, err := http.NewRequest("GET", fmt.Sprintf("/scan/%s", id), nil)

	// inject for testing
	// https://www.gorillatoolkit.org/pkg/mux#SetURLVars
	vars := map[string]string{
		"id": id,
	}
	req = mux.SetURLVars(req, vars)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetScanHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned status code %v want %v", status, http.StatusNotFound)
	}

	expected := fmt.Sprintf("Scan '%s' not found", id)
	if rr.Body.String() != string(expected) {
		t.Errorf("handler returned unexpected body: got %v\n want \n%v",
			rr.Body.String(), expected)
	}
}
