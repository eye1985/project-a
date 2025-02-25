FROM golang:1.24-alpine

WORKDIR /app

COPY . .

EXPOSE 8080

RUN go build -o myapp

CMD ["./myapp"]