BLD_NAME:=mail-sender

.PHONY: build 
build: cmd/sender/main.go 
	go build -o bin/$(BLD_NAME) $<

