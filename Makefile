.PHONY: all

gateway:
	go build -o gateway/cli/gateway.exe gateway/cli/main.go
logic:
	go build -o logic/cli/logic.exe logic/cli/main.go
