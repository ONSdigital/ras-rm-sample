FROM golang:1.14.4-alpine3.12

RUN mkdir "/src"
WORKDIR "/src"

COPY main.go .
COPY sample.go .
COPY go.mod .
COPY go.sum .

RUN go build
RUN ls
CMD "./worker"