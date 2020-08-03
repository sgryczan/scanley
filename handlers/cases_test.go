package handlers

import (
	"github.com/sgryczan/scanley/models"
)

type ScanTest struct {
	Scan     *models.ScanRequest
	Expected string
	Code     int
}

var cases = &[]ScanTest{
	{
		Scan:     &exampleScan,
		Expected: `^{\n.*"scan-id":.*\n}$`,
		Code:     200,
	},
	{
		Scan:     &exampleScan2,
		Expected: `^{\n.*"scan-id":.*\n}$`,
		Code:     200,
	},
	{
		Scan:     &exampleScan3,
		Expected: `^Bad Request: Invalid port range: banana$`,
		Code:     400,
	},
	{
		Scan:     &exampleScan4,
		Expected: `^Bad Request: Port must be between 1-65535$`,
		Code:     400,
	},
}

var exampleScan = models.ScanRequest{
	Hosts: []string{
		"scanme.nmap.org",
		"google.com",
	},
	Ports: []int{
		22,
		80,
		443,
	},
}

var exampleScan2 = models.ScanRequest{
	Hosts: []string{
		"scanme.nmap.org",
		"google.com",
	},
	Ports: []int{
		22,
		80,
		443,
	},
	PortRanges: []string{
		"1-20",
		"200-400",
	},
	Threads: 100,
	Timeout: 3000,
}

var exampleScan3 = models.ScanRequest{
	Hosts: []string{
		"scanme.nmap.org",
		"google.com",
	},
	Ports: []int{
		22,
		80,
		443,
	},
	PortRanges: []string{
		"1-20",
		"banana",
	},
	Threads: 100,
	Timeout: 3000,
}

var exampleScan4 = models.ScanRequest{
	Hosts: []string{
		"scanme.nmap.org",
	},
	Ports: []int{
		22,
		80000,
	},
}
