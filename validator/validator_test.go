package validator_test

import (
	"os"
	"path/filepath"
	"testing"

	"xml-validator/validator"
)

// writeTemp writes content to a temp file and returns its path.
func writeTemp(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

const sampleXSD = `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
  <xs:element name="person">
    <xs:complexType>
      <xs:sequence>
        <xs:element name="name"  type="xs:string"  minOccurs="1" maxOccurs="1"/>
        <xs:element name="age"   type="xs:integer" minOccurs="1" maxOccurs="1"/>
        <xs:element name="email" type="xs:string"  minOccurs="0" maxOccurs="1"/>
      </xs:sequence>
    </xs:complexType>
  </xs:element>
</xs:schema>`

const validXML = `<?xml version="1.0" encoding="UTF-8"?>
<person>
  <name>Jane Doe</name>
  <age>30</age>
  <email>jane@example.com</email>
</person>`

const missingRequiredXML = `<?xml version="1.0" encoding="UTF-8"?>
<person>
  <age>25</age>
</person>`

const wrongTypeXML = `<?xml version="1.0" encoding="UTF-8"?>
<person>
  <name>Bob</name>
  <age>not-a-number</age>
</person>`

const malformedXML = `<?xml version="1.0" encoding="UTF-8"?>
<person>
  <name>Unclosed`

func TestValidateXML_Valid(t *testing.T) {
	dir := t.TempDir()
	xsdPath := writeTemp(t, dir, "schema.xsd", sampleXSD)

	errs, err := validator.ValidateXML([]byte(validXML), xsdPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) != 0 {
		t.Errorf("expected no validation errors, got: %v", errs)
	}
}

func TestValidateXML_MissingRequiredElement(t *testing.T) {
	dir := t.TempDir()
	xsdPath := writeTemp(t, dir, "schema.xsd", sampleXSD)

	errs, err := validator.ValidateXML([]byte(missingRequiredXML), xsdPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) == 0 {
		t.Error("expected validation errors for missing required element, got none")
	}
}

func TestValidateXML_WrongType(t *testing.T) {
	dir := t.TempDir()
	xsdPath := writeTemp(t, dir, "schema.xsd", sampleXSD)

	errs, err := validator.ValidateXML([]byte(wrongTypeXML), xsdPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) == 0 {
		t.Error("expected validation errors for wrong type, got none")
	}
}

func TestValidateXML_MalformedXML(t *testing.T) {
	dir := t.TempDir()
	xsdPath := writeTemp(t, dir, "schema.xsd", sampleXSD)

	errs, err := validator.ValidateXML([]byte(malformedXML), xsdPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errs) == 0 {
		t.Error("expected parse error for malformed XML, got none")
	}
}

func TestValidateXML_BadXSDPath(t *testing.T) {
	_, err := validator.ValidateXML([]byte(validXML), "/nonexistent/schema.xsd")
	if err == nil {
		t.Error("expected error for missing XSD file, got nil")
	}
}
