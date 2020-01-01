FROM golang:1.13.5-buster

COPY . /go/src/github.com/freahs/lunch-server
WORKDIR /go/src/github.com/freahs/lunch-server
RUN go get ./...
CMD ["go", "run",  "main.go", "api.go"]