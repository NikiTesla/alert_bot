FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux go build -o alert_bot  ./cmd/*.go

# Run
CMD ["./alert_bot"]
