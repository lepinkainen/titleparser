FUNCNAME=titleparser
BUILDDIR=build

.PHONY: build

build: test
	env GOOS=linux GOARCH=amd64 go build -o $(BUILDDIR)/$(FUNCNAME)
	cd $(BUILDDIR) && zip $(FUNCNAME).zip $(FUNCNAME)

test:
	go vet ./...
	go test -race -v ./...

lint:
	-golangci-lint run ./...

publish: test lint build
	aws lambda update-function-code --publish --function-name $(FUNCNAME) --zip-file fileb://$(BUILDDIR)/$(FUNCNAME).zip
