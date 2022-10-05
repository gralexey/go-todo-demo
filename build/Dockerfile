FROM golang:1.19.1-bullseye

COPY . /sources
WORKDIR /sources

RUN go build -o main ./main.go

EXPOSE 8081

CMD ["./main"]