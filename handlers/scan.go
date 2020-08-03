package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/xid"
	"github.com/sgryczan/scanley/models"
	"github.com/sgryczan/scanley/scanner"
)

var MaxConcurrentHosts int
var ScanThreads int
var ScanTimeout int

// ScanHandler runs scans
func ScanHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /scan Scan Scan
	//
	// Start a new scan
	// ---
	// consumes:
	// - application/json
	// produces:
	// - text/plain
	// parameters:
	// - name: payload
	//   in: body
	//   description: Scan request. Requires a list of hosts and ports.
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/ScanRequest"
	// responses:
	//   '200':
	//     description: Scan started successfully
	//     type: string
	//   '400':
	//     description: Bad Request
	//     type: string

	scanRequest := &models.ScanRequest{}
	_ = json.NewDecoder(r.Body).Decode(&scanRequest)

	// Expand Port Range
	ports := []int{}
	if len(scanRequest.PortRanges) > 0 {
		pr, err := scanner.ExpandPortRanges(scanRequest.PortRanges)
		if err != nil {
			log.Print(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad Request: %s", err.Error())
			return
		}
		ports = append(ports, pr...)
	}

	ports = append(ports, scanRequest.Ports...)
	ports = scanner.UniquePorts(ports)
	err := scanner.ValidatePorts(ports)

	if scanRequest.Threads == 0 {
		log.Printf("%d", scanRequest.Threads)
		scanRequest.Threads = ScanThreads
	}

	if scanRequest.Timeout == 0 {
		scanRequest.Timeout = ScanTimeout
	}

	if scanRequest.MaxConcurrentHosts == 0 {
		scanRequest.MaxConcurrentHosts = MaxConcurrentHosts
	}

	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request: %s", err.Error())
		return
	}

	id := xid.New()
	resp := &models.NewScanResponse{
		ScanId: id.String(),
	}

	writer := scanner.NewScanWriter()
	log.Printf("Starting scan - id: %s\n", id)

	go scanner.ScanHosts(scanRequest, writer, ports, id.String())

	w.WriteHeader(http.StatusOK)
	json, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintf(w, "%s", json)
}
