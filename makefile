TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build

build: fmt
	./scripts/build.sh

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	gofmt -w $(GOFMT_FILES)

test: fmt
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

acceptance: fmt
	go test -v -i $(TEST) || exit 1
	echo $(TEST) | \
		TF_ACC=true SERVICE_ADDRESS=http://localhost SERVICE_PORT=3001 SERVICE_TOKEN=superSecret xargs -t -n4 go test -v $(TESTARGS) -parallel=4

startapi: fmt
	go run api/main.go