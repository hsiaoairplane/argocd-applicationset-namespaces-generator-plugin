.PHONY: all
all: test build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -o argocd-applicationset-namespaces-generator-plugin -race -v .

.PHONY: clean
clean:
	go clean ./...
