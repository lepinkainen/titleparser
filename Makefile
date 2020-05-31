FUNCNAME=titleparser
BUILDDIR=build

# Include .env and export them as environment variables
-include .env
# Only export the variables we defined in the file
export $(shell sed 's/=.*//' .env)


.PHONY: build

build: test
	env GOOS=linux GOARCH=amd64 go build -o $(BUILDDIR)/$(FUNCNAME)
	cd $(BUILDDIR) && zip $(FUNCNAME).zip $(FUNCNAME)

test:
	go vet ./...
	go test -race -cover -v ./...

lint:
	-golangci-lint run ./...

publish: test lint build
	aws lambda update-function-code --publish --function-name $(FUNCNAME) --zip-file fileb://$(BUILDDIR)/$(FUNCNAME).zip
