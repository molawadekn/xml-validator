FROM golang:1.21-bullseye AS builder

# Install libxml2 dev headers needed by lestrrat-go/libxml2 (cgo)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libxml2-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o xml-validator .

# ── Runtime image ──────────────────────────────────────────────────────────────
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    libxml2 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/xml-validator .
COPY schema.xsd .

EXPOSE 8080
ENV PORT=8080
ENV XSD_PATH=schema.xsd

ENTRYPOINT ["./xml-validator"]
