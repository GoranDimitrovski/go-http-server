FROM golang:1.22

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    && rm -rf /var/lib/apt/lists/*

COPY . .

RUN go build -o app

CMD ["./app"]
