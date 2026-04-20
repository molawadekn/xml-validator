package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"xml-validator/validator"
)

// ValidateResponse is the JSON response returned by the /validate endpoint.
type ValidateResponse struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors,omitempty"`
}

// NewValidateHandler returns an http.HandlerFunc that validates XML against the given XSD path.
func NewValidateHandler(xsdPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, ValidateResponse{
				Valid:  false,
				Errors: []string{"only POST method is allowed"},
			})
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR reading request body: %v", err)
			writeJSON(w, http.StatusBadRequest, ValidateResponse{
				Valid:  false,
				Errors: []string{"failed to read request body"},
			})
			return
		}
		defer r.Body.Close()

		if len(body) == 0 {
			writeJSON(w, http.StatusBadRequest, ValidateResponse{
				Valid:  false,
				Errors: []string{"request body is empty"},
			})
			return
		}

		log.Printf("INFO validating XML document (%d bytes) against %s", len(body), xsdPath)

		errs, err := validator.ValidateXML(body, xsdPath)
		if err != nil {
			log.Printf("ERROR during validation: %v", err)
			writeJSON(w, http.StatusInternalServerError, ValidateResponse{
				Valid:  false,
				Errors: []string{"internal validation error: " + err.Error()},
			})
			return
		}

		if len(errs) > 0 {
			log.Printf("INFO XML is INVALID — %d error(s)", len(errs))
			writeJSON(w, http.StatusBadRequest, ValidateResponse{
				Valid:  false,
				Errors: errs,
			})
			return
		}

		log.Printf("INFO XML is VALID")
		writeJSON(w, http.StatusOK, ValidateResponse{Valid: true})
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("ERROR encoding JSON response: %v", err)
	}
}
