ENDPOINT ?= "http://localhost:8080"

e2e:
	venom run --var url='${ENDPOINT}' venom/*.yml --output-dir test_results