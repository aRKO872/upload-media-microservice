FROM golang:latest

WORKDIR /app/email

COPY . .

RUN go get -d -v ./...

CMD ["go", "run", "main.go"]
