.PHONY: build-ex

build-ex:
	cd ../ && GOOS=linux GOARCH=amd64 go build -o _example/bin/extensions/invocation-history-extension main.go
	chmod +x bin/extensions/invocation-history-extension
	cd bin && zip -r extension.zip extensions/

.PHONY: build-func run
build-func:
	docker build -t invocation-history-ex-func:local .

run:
	docker run \
	--rm \
	-p 9000:8080 \
	invocation-history-ex-func:local