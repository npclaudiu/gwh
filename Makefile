GO ?= go

.PHONY: all
all: gwh

.PHONY: gwh
gwh:
	$(GO) build -o bin/gwh -tags="no_duckdb_arrow" ./cmd/gwh
