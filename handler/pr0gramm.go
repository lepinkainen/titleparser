package handler

import "github.com/lepinkainen/titleparser/lambda"

// Pr0gramm is a weird javascript-only gallery site with no API, just ignore it
func Pr0gramm(url string) (string, error) {
	return "", nil
}

// Register the handler function with corresponding regex
func init() {
	lambda.RegisterHandler(`.*?pr0gramm\.com.*`, Pr0gramm)
}
