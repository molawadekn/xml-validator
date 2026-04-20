# XML Validator — Go REST API

A lightweight HTTP service that validates XML documents against an XSD schema and returns structured JSON responses.

---

## Project Structure

```
xml-validator/
├── main.go                  # Server entry point
├── handler/
│   └── validate.go          # HTTP handler for POST /validate
├── validator/
│   ├── validator.go         # XSD/XML validation logic (pure Go, no HTTP)
│   └── validator_test.go    # Unit tests
├── testdata/
│   ├── valid.xml            # Sample valid XML
│   └── invalid.xml          # Sample invalid XML (triggers errors)
├── schema.xsd               # Default XSD schema (configurable)
├── go.mod
├── Dockerfile
└── README.md
```

---

## Prerequisites

### 1. Install libxml2 (required for CGO binding)

**macOS**
```bash
brew install libxml2
export PKG_CONFIG_PATH="/opt/homebrew/opt/libxml2/lib/pkgconfig"
```

**Ubuntu / Debian**
```bash
sudo apt-get install -y libxml2-dev pkg-config
```

**Alpine Linux**
```bash
apk add libxml2-dev pkgconfig
```

### 2. Install Go ≥ 1.21
Download from https://go.dev/dl/

---

## Getting Started

```bash
# Clone / enter the project
cd xml-validator

# Download dependencies
go mod tidy

# Run the server (uses schema.xsd in current directory by default)
go run main.go
```

The server starts on **http://localhost:8080**.

### Configuration via environment variables

| Variable   | Default      | Description                          |
|------------|--------------|--------------------------------------|
| `PORT`     | `8080`       | TCP port to listen on                |
| `XSD_PATH` | `schema.xsd` | Path to the XSD schema file          |

```bash
PORT=9090 XSD_PATH=/etc/schemas/myschema.xsd go run main.go
```

---

## API Reference

### `POST /validate`

**Request**
- Content-Type: `application/xml` (or any; body must be raw XML)
- Body: XML document to validate

**Responses**

| Status | Meaning                              |
|--------|--------------------------------------|
| `200`  | XML is valid                         |
| `400`  | XML is invalid or malformed          |
| `405`  | Wrong HTTP method                    |
| `500`  | Internal error (e.g. bad XSD file)   |

**200 — Valid**
```json
{ "valid": true }
```

**400 — Invalid**
```json
{
  "valid": false,
  "errors": [
    "Element 'age': 'not-a-number' is not a valid value of the atomic type 'xs:integer'.",
    "Element 'person': Missing child element(s). Expected is ( name )."
  ]
}
```

### `GET /health`

Returns `{"status":"ok"}` — useful for Docker / Kubernetes health checks.

---

## curl Examples

### Valid XML (expect HTTP 200)
```bash
curl -s -X POST http://localhost:8080/validate \
  -H "Content-Type: application/xml" \
  --data-binary @testdata/valid.xml | jq .
```
**Response:**
```json
{ "valid": true }
```

---

### Invalid XML — missing required element and wrong type (expect HTTP 400)
```bash
curl -s -X POST http://localhost:8080/validate \
  -H "Content-Type: application/xml" \
  --data-binary @testdata/invalid.xml | jq .
```
**Response:**
```json
{
  "valid": false,
  "errors": [
    "Element 'age': 'not-a-number' is not a valid value of the atomic type 'xs:integer'.",
    "Element 'person': Missing child element(s). Expected is ( name )."
  ]
}
```

---

### Inline XML string — valid
```bash
curl -s -X POST http://localhost:8080/validate \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?><person><name>Alice</name><age>28</age></person>' | jq .
```

---

### Inline XML string — invalid (age is not an integer)
```bash
curl -s -X POST http://localhost:8080/validate \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?><person><name>Bob</name><age>twenty</age></person>' | jq .
```

---

### Malformed XML (expect HTTP 400 with parse error)
```bash
curl -s -X POST http://localhost:8080/validate \
  -H "Content-Type: application/xml" \
  -d '<?xml version="1.0"?><person><name>Unclosed' | jq .
```

---

### Empty body (expect HTTP 400)
```bash
curl -s -X POST http://localhost:8080/validate | jq .
```

---

## Running Tests

```bash
go test ./...
```

Or with verbose output:
```bash
go test -v ./validator/...
```

---

## Docker

### Build
```bash
docker build -t xml-validator .
```

### Run
```bash
docker run -p 8080:8080 xml-validator
```

### Run with a custom XSD
```bash
docker run -p 8080:8080 \
  -v /path/to/your/schema.xsd:/app/schema.xsd \
  -e XSD_PATH=/app/schema.xsd \
  xml-validator
```

---

## Sample XSD (`schema.xsd`)

The default schema expects a `<person>` document:

```xml
<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
  <xs:element name="person">
    <xs:complexType>
      <xs:sequence>
        <xs:element name="name"  type="xs:string"  minOccurs="1"/>
        <xs:element name="age"   type="xs:integer" minOccurs="1"/>
        <xs:element name="email" type="xs:string"  minOccurs="0"/>
      </xs:sequence>
    </xs:complexType>
  </xs:element>
</xs:schema>
```

Swap in your own XSD by setting `XSD_PATH` — no code changes needed.
