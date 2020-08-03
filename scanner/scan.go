package scanner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/sgryczan/scanley/models"
)

var fileMutex sync.Mutex

type Scanner struct {
	ConcurrentHosts int
	ThreadsPerHost  int
}

type scanWriter struct {
	mutex *sync.Mutex
}

func NewScanWriter() *scanWriter {
	return &scanWriter{mutex: &sync.Mutex{}}
}

func worker(d *net.Dialer, host string, ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", host, p)
		conn, err := d.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func hostWorker(w *scanWriter, r *models.ScanRequest, ports []int, hosts, results chan string, id string) {
	for h := range hosts {
		Scan(w, r.Threads, r.Timeout, h, ports, id)
		results <- h
	}
}

//ScanHosts handles concurrent scans for a group of hosts
func ScanHosts(r *models.ScanRequest, w *scanWriter, ports []int, id string) {
	hosts := r.Hosts
	threads := min(len(r.Hosts), r.MaxConcurrentHosts)
	results := make(chan string)
	hostBuffer := make(chan string, threads)
	scannedHosts := []string{}

	log.Printf("Processing Scan Request %s - %d hosts on %d thread(s)\n", id, len(r.Hosts), threads)

	for i := 0; i < cap(hostBuffer); i++ {
		go hostWorker(w, r, ports, hostBuffer, results, id)
	}

	go func() {
		for _, i := range hosts {
			hostBuffer <- i
		}
	}()

	for _, _ = range ports {
		host := <-results
		scannedHosts = append(scannedHosts, host)
	}

	close(hostBuffer)
	close(results)

	return
}

// Scan runs a port scan against a remote host
func Scan(w *scanWriter, num_threads int, timeout int, host string, ports []int, id string) {
	threads := min(len(ports), num_threads)
	portBuffer := make(chan int, threads)
	results := make(chan int)
	dialer := net.Dialer{Timeout: time.Millisecond * time.Duration(timeout)}
	log.Printf("Scan - id: %s host: %s ports: %d threads: %d\n", id, host, len(ports), threads)
	var openports []int

	for i := 0; i < cap(portBuffer); i++ {
		go worker(&dialer, host, portBuffer, results)
	}

	go func() {
		for _, i := range ports {
			portBuffer <- i
		}
	}()

	for _, _ = range ports {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(portBuffer)
	close(results)

	sort.Ints(openports)
	scanResults := &models.ScanResults{
		Date: time.Now().Format(time.RFC3339),
		Results: []models.HostScan{
			{
				Host:  host,
				Ports: openports,
			},
		},
	}
	err := w.Update(scanResults, id)
	if err != nil {
		return
	}

	log.Printf("Scan complete - id: %s, host: %s\n", id, host)
	return
}

// Update writes scan results for a single host to a file. If the file exists, it
// will be appended to.
func (w *scanWriter) Update(scan *models.ScanResults, filename string) error {
	w.mutex.Lock()
	if _, err := os.Stat("scans/" + filename); err == nil {
		data, err := ReadScanFromFile("scans/" + filename)
		if err != nil {
			log.Print(err.Error())
			return err
		}
		scan.Results = append(scan.Results, data.Results...)
	}
	contents, err := json.MarshalIndent(scan, "", "  ")
	if err != nil {
		log.Print(err.Error())
		return err
	}
	err = ioutil.WriteFile("scans/"+filename, contents, 0644)
	w.mutex.Unlock()
	return err
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
