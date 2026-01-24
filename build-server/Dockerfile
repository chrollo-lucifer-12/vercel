FROM ubuntu:focal


RUN apt-get update && \
    apt-get install -y curl git ca-certificates && \
    rm -rf /var/lib/apt/lists/*


RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs

ENV GO_VERSION=1.22.1
RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz -o go.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /home/app

COPY go.mod go.sum ./
RUN go mod download

COPY src/ src/
RUN go build -o app ./src

COPY main.sh main.sh
RUN chmod +x main.sh

ENTRYPOINT ["/home/app/main.sh"]
