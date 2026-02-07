FROM golang:1.25-alpine AS builder

WORKDIR /home/app

RUN apk add --no-cache git

COPY go.work.docker ./go.work
COPY build-server/go.mod build-server/go.sum ./build-server/
COPY shared/go.mod shared/go.sum ./shared/


RUN go work sync
RUN cd build-server && go mod download

COPY build-server ./build-server
COPY shared ./shared


WORKDIR /home/app/build-server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bs .



FROM alpine:3.19

WORKDIR /

RUN apk add --no-cache ca-certificates git nodejs npm dos2unix


RUN mkdir -p /code /home/app/output


COPY --from=builder /bs /bs
COPY build-server/main.sh /main.sh

RUN dos2unix /main.sh
RUN chmod +x /main.sh

ENTRYPOINT ["/main.sh"]
