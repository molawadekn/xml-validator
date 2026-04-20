package validator

import (
	"fmt"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
)

// ValidateXML validates xmlData against the XSD schema at xsdPath.
// Returns a slice of human-readable error strings, or an error if the
// validation machinery itself fails (e.g. cannot parse the schema).
func ValidateXML(xmlData []byte, xsdPath string) ([]string, error) {
	// Parse the XSD schema
	schema, err := xsd.ParseFromFile(xsdPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XSD schema: %w", err)
	}
	defer schema.Free()

	// Parse the XML document
	doc, err := libxml2.ParseString(string(xmlData))
	if err != nil {
		// ParseString returns a (doc, error) where the error may contain
		// parse-level errors. We surface them as validation errors.
		return []string{fmt.Sprintf("XML parse error: %s", err.Error())}, nil
	}
	defer doc.Free()

	// Validate the parsed document against the schema
	if err := schema.Validate(doc); err != nil {
		return extractErrors(err), nil
	}

	return nil, nil
}

// extractErrors unpacks a libxml2 validation error into individual strings.
func extractErrors(err error) []string {
	type multiError interface {
		Errors() []error
	}

	if me, ok := err.(multiError); ok {
		errs := me.Errors()
		messages := make([]string, 0, len(errs))
		for _, e := range errs {
			messages = append(messages, e.Error())
		}
		return messages
	}

	return []string{err.Error()}
}
