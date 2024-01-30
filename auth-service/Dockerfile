FROM golang:latest

WORKDIR /app/auth

COPY . .

RUN go get -d -v ./...

CMD ["go", "run", "main.go"]