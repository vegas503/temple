BINDIR = ./bin
TARGET = $(BINDIR)/temple
SRC = $(shell find . -name *.go)

$(TARGET): $(SRC) go.mod
	@mkdir -p $(BINDIR)
	go build -trimpath -ldflags '-s -w -buildid=' -gcflags=all="-B -l -wb=false" -o $@ ./cmd/temple
	# tinygo build -o $@ ./cmd/main

.PHONY: clean
clean:
	rm -rf $(BINDIR)

.PHONY: lint
lint:
	golangci-lint run
