// passthrough.go
package main

import (
	"log"
	"net/http"
	"net/http/httputil" // Required for the ReverseProxy
	"net/url"
	"time"
)

// The target URL for the passthrough.
// For this example, we'll use a public echo service.
const targetURL = "http://3.82.251.214:8688/service/EchoService"

func main() {
	// Parse the target URL.
	remote, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	// Create a new ReverseProxy.
	// This handles the request forwarding for us.
	proxy := httputil.NewSingleHostReverseProxy(remote)

	// Custom transport to log requests and handle errors.
	proxy.Transport = &http.Transport{
		// Set a timeout for the client request.
		// It's good practice to have a timeout to prevent hanging connections.
		// The Ballerina service has this implicitly.
		IdleConnTimeout: 10 * time.Second,
		MaxIdleConns:    100,
	}

	// Define a custom handler function for our proxy.
	// This is where we can add custom logic, like error handling.
	http.HandleFunc("/passthrough", func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request.
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		
		// The http.Request body is a one-time stream.
		// httputil.NewSingleHostReverseProxy will handle reading the body and
		// streaming it to the target, so we don't need to read it ourselves.
		
		// Forward the request to the target.
		// httputil.NewSingleHostReverseProxy.ServeHTTP handles all the details:
		// - It copies the incoming request's method, headers, and body.
		// - It sends the request to the target.
		// - It copies the response from the target back to the original client.
		proxy.ServeHTTP(w, r)
	})

	// Start the HTTP server.
	// We'll use log.Fatal to handle any errors from ListenAndServe.
	log.Println("Starting passthrough service on :9090")
	if err := http.ListenAndServe(":8688", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

