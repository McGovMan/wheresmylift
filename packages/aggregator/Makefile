current_branch := $(shell git branch --show-current)-$(shell git rev-parse --short HEAD)
untracked_files := $(shell git ls-files . --exclude-standard --others)
untracked_files_length := $(shell echo -n '$(untracked_files)' | wc -c)

units:
	go test -timeout 10s -p 1 -v -count=1 -failfast ./...

coverage:
	go test -timeout 10s -p 1 -v -cover -count=1 -failfast $(shell go list ./... | grep -v -E 'aggregator$$') -coverprofile cover.out
	go tool cover -func cover.out | grep total | grep -Eo '[0-9]+\.[0-9]+'

lint:
	golangci-lint run -v
