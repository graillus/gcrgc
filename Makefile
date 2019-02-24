build:
	go build -o bin/gcrgc cmd/gcrgc/*

test:
	go test -v -cover cmd/gcrgc/*
