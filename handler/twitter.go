package handler

import "github.com/lepinkainen/titleparser/lambda"

// ApinaBiz titles are always useless, just don't return anything
func Twitter(url string) (string, error) {
	return "", nil
}

// Register the handler function with corresponding regex
func init() {
	lambda.RegisterHandler(".*?twitter.com.*", Twitter)
}
