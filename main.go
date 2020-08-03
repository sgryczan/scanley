package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"github.com/sgryczan/scanley/handlers"
)

var (
	listenPort         = flag.Int("port", 8080, "Port to listen on")
	threads            = flag.Int("threads", 500, "Number of threads per host scan")
	timeout            = flag.Int("timeout", 3000, "Timeout duration (milliseconds)")
	maxConcurrentHosts = flag.Int("max-concurrent-hosts", 10, "Number of hosts scans to run concurrently per request")
)

func genXid() {
	id := xid.New()
	fmt.Printf("github.com/rs/xid:              %s\n", id.String())
}

func main() {
	flag.Parse()
	version := handlers.Version
	handlers.ScanThreads = *threads
	handlers.ScanTimeout = *timeout
	handlers.MaxConcurrentHosts = *maxConcurrentHosts

	r := mux.NewRouter()

	fmt.Printf("Started scanley v%s\n", version)
	fmt.Printf("Threads per scan: %d\n", handlers.ScanThreads)
	fmt.Printf("Timeout: %d\n", *timeout)

	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("/about", handlers.AboutHandler)
	r.HandleFunc("/scan", handlers.ListHandler).Methods("GET")
	r.HandleFunc("/scan", handlers.ScanHandler).Methods("POST")
	r.HandleFunc("/scan/{id}", handlers.GetScanHandler).Methods("GET")

	sh := http.StripPrefix("/api",
		http.FileServer(http.Dir("./swaggerui/")))
	r.PathPrefix("/api/").Handler(sh)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + strconv.Itoa(*listenPort),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
