package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sgryczan/scanley/models"
	"github.com/sgryczan/scanley/scanner"
)

func TestListHandler(t *testing.T) {
	// set up a temp inventory dir
	targetDir := "scans/"
	err := os.Mkdir(targetDir, 0777)
	defer os.RemoveAll(targetDir)

	res := &models.ScanResults{
		Date: "now",
	}
	writer := scanner.NewScanWriter()

	for i := 1; i < 6; i++ {
		id := fmt.Sprintf("testScan%d", i)
		err = writer.Update(res, id)
	}

	if err != nil {
		t.Fatalf("%v", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("/scan"), nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned status code %v want %v", status, http.StatusOK)
	}

	expected := `{
  "count": 5,
  "scans": [
    "testScan1",
    "testScan2",
    "testScan3",
    "testScan4",
    "testScan5"
  ]
}`

	if rr.Body.String() != string(expected) {
		t.Errorf("handler returned unexpected body: got \n%v\n want \n%v\n",
			rr.Body.String(), expected)
	}
}
