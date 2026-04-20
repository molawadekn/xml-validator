package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"xml-validator/handler"
)

func main() {
	port := getEnv("PORT", "8080")
	xsdPath := getEnv("XSD_PATH", "schema.xsd")

	// Verify XSD file exists at startup
	if _, err := os.Stat(xsdPath); os.IsNotExist(err) {
		log.Fatalf("XSD schema file not found at path: %s", xsdPath)
	}

	log.Printf("Starting XML Validator API on port %s", port)
	log.Printf("Using XSD schema: %s", xsdPath)

	mux := http.NewServeMux()
	mux.HandleFunc("/validate", handler.NewValidateHandler(xsdPath))
	mux.HandleFunc("/health", healthHandler)

	addr := fmt.Sprintf(":%s", port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
