.PHONY: build test

build:
	@go build

test:
	go test -v ./...

run: build
	@chmod +x ./leaguerank
	@./leaguerank