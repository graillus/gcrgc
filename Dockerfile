FROM golang:alpine as build

RUN apk add make

WORKDIR /go/src/github.com/graillus/gcrgc
COPY . .

RUN make build

FROM alpine

COPY --from=build /go/src/github.com/graillus/gcrgc/bin/gcrgc /usr/local/bin/gcrgc
