package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sgryczan/scanley/models"
	"github.com/sgryczan/scanley/scanner"
)

// ListHandler lists all scans
func ListHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /scan Scan Scan
	//
	// Lists all scans
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     description: Success
	//     type: string
	inv := models.ScanInventory{}

	scans, err := scanner.ListScanFiles()
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	inv.Count = len(scans)
	inv.Scans = scans

	res, err := json.MarshalIndent(inv, "", "  ")

	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", res)
}
