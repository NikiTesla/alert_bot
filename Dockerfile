FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o alert_bot  ./cmd/*.go

EXPOSE 2704

# Run
CMD ["./alert_bot"]