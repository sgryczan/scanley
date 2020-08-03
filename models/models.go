package models

// Scan represents a bootable host
// swagger:model
type ScanResults struct {
	Date    string     `json:"date"`
	Results []HostScan `json:"results"`
}

type HostScan struct {
	Host  string `json:"host"`
	Ports []int  `json:"openports"`
}

// ScanRequest represent the input for a scan request
// swagger:model
type ScanRequest struct {
	// required: true
	// example: ["scanme.nmap.org", "google.com", "amazon.com"]
	Hosts []string `json:"hosts"`
	// required: true
	// example: [25, 80, 443]
	Ports []int `json:"ports"`
	// required: false
	// example: ["1-1024"]
	PortRanges []string `json:"port-ranges"`
	// required: false
	// example: 500
	Threads int `json:"threads"`
	// required: false
	// example: 3000
	Timeout int `json:"connect-timeout"`
	// required: false
	// example: 10
	MaxConcurrentHosts int `json:"max-concurrent-hosts"`
}

// NewScanResponse is the response to a new scan
// swagger:model
type NewScanResponse struct {
	ScanId string `json:"scan-id"`
}

// ScanInventory is a collection of scans
// swagger:model
type ScanInventory struct {
	Count int      `json:"count"`
	Scans []string `json:"scans"`
}
