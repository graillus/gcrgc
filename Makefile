build:
	go build -o bin/gcrgc cmd/gcrgc/gcrgc.go

test:
	go test -v -cover cmd/gcrgc/*
