FROM ubuntu:focal


RUN apt-get update && \
    apt-get install -y curl git ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ENV GO_VERSION=1.22.1

RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz -o go.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz


ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH=/go
ENV PATH="$GOPATH/bin:${PATH}"

WORKDIR /home/app


COPY main.sh main.sh
COPY . .

RUN chmod +x main.sh

RUN go mod download

ENTRYPOINT ["/home/app/main.sh"]
