FROM golang:1.22 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
# Copy the go source
COPY api/ api/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o api-server api/main.go

# Use distroless as minimal base image to package the api server binary
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/api-server .
USER 65532:65532

ENTRYPOINT ["/api-server"]
