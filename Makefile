BINDIR = ./bin
TARGET = $(BINDIR)/temple
SRC = $(shell find . -name *.go)

$(TARGET): $(SRC) go.mod
	@mkdir -p $(BINDIR)
	go build -trimpath -ldflags '-s -w -buildid=' -gcflags=all="-B -l" -o $@ ./cmd/temple
	# tinygo build -o $@ ./cmd/main

clean:
	rm -rf $(BINDIR)
