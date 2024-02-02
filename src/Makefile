DESTINATION_BUCKET ?= "bkt-prj-eng-d-adform-sync-svc-ewh5-int-test-dest"
SOURCE_BUCKET ?= "bkt-prj-eng-d-adform-sync-svc-ewh5-int-test-src"
OBJECT ?= ""
PROJECT_ID ?= "prj-eng-d-adform-sync-svc-ewh5"

build: *.go fmt
	@go build -o build/bin/redsyncungzip ./cmd/main.go

fmt:
	@gofmt -w=true -s $$(find . -type f -name '*.go')
	@goimports -w=true -d $$(find . -type f -name '*.go')

lint:
	@echo "Linting Go files..."
	@golangci-lint run --deadline 5m 
	@echo "Linting complete"

setup:
	go install golang.org/x/tools/cmd/goimports@latest

test:
	@echo "Start Unit Testing..."
	@go test -v

bench:
	@go test -bench=.

ci: build test
	@echo "Start Integration Testing..."
	@chmod +x build/bin/redsyncungzip
	@./build/bin/redsyncungzip -p $(PROJECT_ID) -d $(DESTINATION_BUCKET) -s $(SOURCE_BUCKET) -o $(OBJECT)
