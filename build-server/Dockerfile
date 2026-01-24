FROM golang:1.25-alpine AS builder

WORKDIR /home/app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY src/ src/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bs ./src



FROM alpine:3.19

WORKDIR /

RUN apk add --no-cache ca-certificates git nodejs npm dos2unix

# Create the folders
RUN mkdir -p /code /home/app/output

# Copy the binary to /code
COPY --from=builder /bs /bs

# Copy main.sh
COPY main.sh /main.sh

RUN dos2unix /main.sh
RUN chmod +x /main.sh

ENTRYPOINT ["/main.sh"]
