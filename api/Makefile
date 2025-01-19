current_branch := $(shell git branch --show-current)-$(shell git rev-parse --short HEAD)
untracked_files := $(shell git ls-files . --exclude-standard --others)
untracked_files_length := $(shell echo -n '$(untracked_files)' | wc -c)

test:
	go test -timeout 10s -p 1 -v -count=1 -failfast ./...

coverage:
	go test -timeout 10s -p 1 -v -cover -count=1 -failfast $(shell go list ./... | grep -v -E 'api$$' | grep -v -E 'utils$$') -coverprofile cover.out

coverage-total:
	go tool cover -func cover.out | grep total | grep -Eo '[0-9]+\.[0-9]+'

lint:
	golangci-lint run -v

generate-swagger:
	go get . || true
	go install github.com/swaggo/swag/cmd/swag@v1.16.4
	swag init --dir ".,internal/server,internal/helpers"

verify-swagger:
	# Using /bin/cp as cp seems to be aliased and is losing flags on the way
	/bin/cp -Rf docs /tmp/docs_branch
	cmp --silent docs/docs.go /tmp/docs_branch/docs.go || (echo "Swagger CLI was not ran" && exit 1)
	cmp --silent docs/swagger.json /tmp/docs_branch/swagger.json || (echo "Swagger CLI was not ran" && exit 1)
	cmp --silent docs/swagger.yaml /tmp/docs_branch/swagger.yaml || (echo "Swagger CLI was not ran" && exit 1)
	echo "Swagger has been verified successfully"
