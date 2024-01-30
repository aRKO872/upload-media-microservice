FROM golang:latest

WORKDIR /app/db

COPY . .

RUN go get -d -v ./...

CMD ["go", "run", "main.go"]
