# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . .

# Build
RUN CGO_ENABLED=0 go build -a -o manager main.go

FROM alpine:latest

RUN apk update && apk upgrade && apk add git openssh-client && ssh-keyscan github.ibm.com > /etc/ssh/ssh_known_hosts

WORKDIR /
COPY --from=builder /workspace/manager .

RUN adduser --disabled-password -u 65532 main 65532
USER 65532:65532

ENTRYPOINT ["/manager"]
