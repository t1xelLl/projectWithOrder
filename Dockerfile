FROM golang:1.24.1-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o projectWithOrder ./cmd/main.go

COPY web ./web

EXPOSE 8080

CMD ["./projectWithOrder"]