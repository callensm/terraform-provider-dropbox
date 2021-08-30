clean:
	rm -rf vendor/

test: vendor
	go test -v ./...

vendor: clean
	go mod tidy && go mod vendor

.PHONY: clean test vendor
