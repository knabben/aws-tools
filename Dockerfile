FROM golang:1.7.4-alpine

RUN apk update
RUN apk add curl git 

WORKDIR /app
ADD . /app

RUN go get -u -v github.com/knabben/aws-tools
RUN go build main.go

ENTRYPOINT ["./main"]
