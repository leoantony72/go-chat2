FROM golang:alpine

RUN mkdir /app

WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN go mod download
ADD . .

EXPOSE 6300
CMD go run main.go
