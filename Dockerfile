# syntax=docker/dockerfile:1

# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod ./
COPY go.sum ./
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY *.go ./

# Build
RUN go build -a -o udp-server 

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/udp-server .
USER nonroot:nonroot
#EXPOSE 8080

ENTRYPOINT ["/udp-server"]
