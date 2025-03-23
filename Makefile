BINDIR = ./bin
TARGET = $(BINDIR)/temple
SRC = $(shell find . -name *.go)

$(TARGET): $(SRC) go.mod
	@mkdir -p $(BINDIR)
	CGO_ENABLED=0 go build -v -trimpath -ldflags '-s -w -buildid=' -gcflags=all="-B -l -wb=false" -o $@ ./cmd/temple
	# tinygo build -o $@ ./cmd/main

.PHONY: clean
clean:
	rm -rf $(BINDIR)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	cat examples/template.tpl | USER=foo COLORS=red,green,blue FEATURES=one=y,two=n,three=y go run ./cmd/temple | diff - ./examples/expected.txt
