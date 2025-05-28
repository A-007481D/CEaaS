# Build the controller binary
FROM golang:1.22 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
# Copy the go source
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o controller cmd/controller/main.go

# Use distroless as minimal base image to package the controller binary
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/controller .
USER 65532:65532

ENTRYPOINT ["/controller"]
