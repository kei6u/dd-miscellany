FROM golang:1.17 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY main.go main.go

WORKDIR /workspace
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/app .

# Gnerate a Self-Signed Certificate by following command.
# openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout go-server.key -out go-server.crt
COPY go-server.crt go-server.crt
COPY go-server.key go-server.key

USER 65532:65532

ENTRYPOINT ["/app"]
