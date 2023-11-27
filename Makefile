FUNCNAME=titleparser
BINARYNAME=bootstrap
BUILDDIR=build

# Include .env and export them as environment variables
-include .env
# Only export the variables we defined in the file
export $(shell touch .env && sed 's/=.*//' .env)


.PHONY: build

build_local: test
	go build -o $(FUNCNAME)

build: test
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(BUILDDIR)/$(BINARYNAME)
	cd $(BUILDDIR) && zip $(FUNCNAME).zip $(BINARYNAME)

test:
	go vet ./...
	go test -cover -v ./...

lint:
	-golangci-lint run ./...

publish: test lint build
	aws lambda update-function-code --publish --function-name $(FUNCNAME) --zip-file fileb://$(BUILDDIR)/$(FUNCNAME).zip
