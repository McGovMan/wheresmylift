current_branch := $(shell git branch --show-current)-$(shell git rev-parse --short HEAD)
untracked_files := $(shell git ls-files . --exclude-standard --others)
untracked_files_length := $(shell echo -n '$(untracked_files)' | wc -c)

ENDPOINT ?= "http://localhost:8080"

units:
	go test -timeout 10s -p 1 -v -count=1 -failfast ./...

coverage:
	go test -timeout 10s -p 1 -v -cover -count=1 -failfast $(shell go list ./... | grep -v -E 'api$$' | grep -v -E 'utils|docs$$') -coverprofile cover.out
	go tool cover -func cover.out | grep total | grep -Eo '[0-9]+\.[0-9]+'

lint:
	golangci-lint run -v

e2e:
	venom run --var url='${ENDPOINT}' venom/*.yml --output-dir test_results

generate-swagger:
	rm -R docs || true
	swag init --dir ".,internal/server,internal/helpers"

verify-swagger:
	rm -R /tmp/docs_branch || true
	# Using /bin/cp as cp seems to be aliased and is losing flags on the way
	/bin/cp -Rf docs /tmp/docs_branch
	$(MAKE) generate-swagger
	cmp --silent docs/docs.go /tmp/docs_branch/docs.go || (echo "Swagger CLI was not ran" && exit 1)
	cmp --silent docs/swagger.json /tmp/docs_branch/swagger.json || (echo "Swagger CLI was not ran" && exit 1)
	cmp --silent docs/swagger.yaml /tmp/docs_branch/swagger.yaml || (echo "Swagger CLI was not ran" && exit 1)
	echo "Swagger has been verified successfully"
