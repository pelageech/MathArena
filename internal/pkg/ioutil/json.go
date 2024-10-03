package ioutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/log"

	"github.com/pelageech/matharena/internal/models"
)

// JSONErrorWriter interpreters an error in json format
// and write to http.ResponseWriter.
type JSONErrorWriter struct {
	Logger *log.Logger
}

// Error is a helper function to write a generic error
// to the response writer with the given status code and message.
func (ew JSONErrorWriter) Error(w http.ResponseWriter, message string, status int) {
	// write http status code
	w.WriteHeader(status)

	// write json response
	// strings never returns error on marshal
	err := ToJSON(models.GenericError{
		Message: message,
	}, w)
	if err != nil {
		ew.Logger.Error("Unable to write JSON response", "error", err)
	}
}

// ToJSON serializes the given interface into a string-based JSON format.
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)

	err := e.Encode(i)
	if err != nil {
		return fmt.Errorf("unable to encode JSON: %w", err)
	}

	return nil
}

// FromJSON deserializes the object from JSON string
// in an io.Reader to the given interface.
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)

	err := d.Decode(i)
	if err != nil {
		return fmt.Errorf("unable to decode JSON: %w", err)
	}

	return nil
}
