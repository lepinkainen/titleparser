FUNCNAME=titleparser

build: clean
	env GOOS=linux GOARCH=amd64 go build -o $(FUNCNAME)
	zip $(FUNCNAME).zip $(FUNCNAME)

clean:
	rm -f $(FUNCNAME).zip

test:
	go test -v

publish: test build
	aws lambda update-function-code --function-name $(FUNCNAME) --zip-file fileb://$(FUNCNAME).zip
