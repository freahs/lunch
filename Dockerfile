FROM golang:1.13.5-buster

COPY . /go/src/github.com/freahs/lunch-server
RUN cd /go/src/github.com/freahs/lunch-server && go install