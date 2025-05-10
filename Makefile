.PHONY: build test

build:
	@go build

test:
	go test -v ./..

run: build
	@chmod +x ./league_rank
	@./league_rank