package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sgryczan/scanley/scanner"
)

// GetScanHandler collects scans
func GetScanHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /scan/{id} Scan Scan
	//
	// Retrieves a scan
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: ID of scan to be returned.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: Scan details retrieved successfully
	//     type: string
	// responses:
	//   '200':
	//     description: Scan retrieved successfully
	//     type: string
	//   '400':
	//     description: Scan configuration not found
	//     type: string
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("GET /scan/%s", id)
	scanRequest, err := scanner.ReadScanFromFile("scans/" + id)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Scan '%s' not found", id)
		return
	}

	res, err := json.MarshalIndent(scanRequest, "", "  ")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", res)
	log.Print(fmt.Sprintf("Retrieved scan id: %s", id))
}
