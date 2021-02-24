
.PHONY: test
test:
	@go test -race -count=1 -cover $$(go list ./... | grep -Ev 'mocks')

.PHONY: build
build:
	@go build ./...
	@cd cmd/azul-cli && go build

.PHONY: clean
clean:
	@go clean ./...

.PHONY: tidy
tidy:
	@go mod tidy