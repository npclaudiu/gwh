GO ?= go

.PHONY: all
all: gwh

.PHONY: gwh
gwh:
	$(GO) build -o bin/gwh ./cmd/gwh
