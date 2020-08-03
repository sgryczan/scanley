package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
)

func TestScanHandler(t *testing.T) {
	targetDir := "scans/"
	_ = os.Mkdir(targetDir, 0777)
	defer os.RemoveAll(targetDir)

	for _, c := range *cases {
		rawJSON, err := json.Marshal(c.Scan)
		if err != nil {
			t.Fatalf("error marshalling json: %v", err)
		}

		req, err := http.NewRequest("POST", "/scan", bytes.NewBuffer(rawJSON))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ScanHandler)

		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != c.Code {
			t.Errorf("handler returned status code %v want %v", status, c.Code)
		}

		match, _ := regexp.MatchString(c.Expected, rr.Body.String())
		if !(match) {
			t.Errorf("handler returned unexpected body: got %v\n want pattern match:\n%v",
				rr.Body.String(), c.Expected)
		}
	}
}
