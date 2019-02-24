FROM golang:alpine as build

WORKDIR /go/src/github.com/graillus/gcrgc
COPY . .

RUN go build -v -o bin/gcrgc *.go

FROM google/cloud-sdk:alpine

COPY --from=build /go/src/github.com/graillus/gcrgc/bin/gcrgc /usr/local/bin/gcrgc
