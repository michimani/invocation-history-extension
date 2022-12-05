.PHONY: build

build:
	GOOS=linux GOARCH=amd64 go build -o bin/extensions/invocation-history-extension main.go
	chmod +x bin/extensions/invocation-history-extension
	cd bin && zip -r extension.zip extensions/
