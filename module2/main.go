package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	// Enable profiling, visit endpoint /debug/pprof for stats
	_ "net/http/pprof"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Copy all headers from request to response
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("Header: %s, Value: %s\n", name, value)
			w.Header().Add(name, value)
		}
	}

	// Add VERSION from environment variables to response header
	version := os.Getenv("VERSION")
	if version != "" {
		w.Header().Add("VERSION", version)
	}

	// Log client IP
	clientIP := getRemoteIp(r)
	fmt.Printf("Client IP: %s", clientIP)
}

func getRemoteIp(r *http.Request) string {
	// reference: https://stackoverflow.com/a/55738279
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
		IPAddresses := strings.Split(IPAddress, ",")

		// Return the first IP address which is the client IP
		return strings.TrimSpace(IPAddresses[0])
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/healthz", healthHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Start httpserver failed, error: %v", err.Error())
	}
}
