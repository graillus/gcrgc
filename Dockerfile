FROM golang:alpine as build

WORKDIR /go/src/github.com/graillus/gcrgc
COPY . .

RUN go build -o bin/gcrgc cmd/gcrgc/gcrgc.go

FROM alpine

COPY --from=build /go/src/github.com/graillus/gcrgc/bin/gcrgc /usr/bin/gcrgc

ENTRYPOINT ["/usr/bin/gcrgc"]
