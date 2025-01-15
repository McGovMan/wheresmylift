current_branch := $(shell git branch --show-current)-$(shell git rev-parse --short HEAD)
untracked_files := $(shell git ls-files . --exclude-standard --others)
untracked_files_length := $(shell echo -n '$(untracked_files)' | wc -c)

test:
	go test -timeout 10s -p 1 -v -count=1 -failfast ./...

coverage:
	go test -timeout 10s -p 1 -v -cover -count=1 -failfast $(shell go list ./... | grep -v -E 'api$$' | grep -v -E 'utils$$') -coverprofile cover.out

coverage-total:
	go tool cover -func cover.out | grep total | grep -Eo '[0-9]+\.[0-9]+'

build:
ifneq "$(untracked_files_length)" "0"
	@echo "All files should be commited before building"
	exit 1
endif

	go build -ldflags "-X github.com/mcgovman/wheresmylift/api.Version=$(current_branch)"

lint:
	golangci-lint run -v
