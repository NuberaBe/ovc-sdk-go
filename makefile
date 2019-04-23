TEST?=./ovc/...

default: test

generate:
	go generate $(TEST)

test: lint generate
	go test $(TEST) -timeout=30s -parallel=4

lint: fmtcheck
	@echo "==> Checking source code against linters..."
	@golangci-lint run $(TEST)
	@go vet $(TEST)

fmtcheck:
	@echo "==> Checking that code complies with gofmt requirements...."
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

tools:
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: default test lint fmtcheck fmt tools
