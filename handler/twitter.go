package handler

import "github.com/lepinkainen/titleparser/lambda"

// Twitter is blocking external agents, so just return empty
func Twitter(url string) (string, error) {
	return "", nil
}

// Register the handler function with corresponding regex
func init() {
	lambda.RegisterHandler(`.*?twitter\.com.*`, Twitter)
}
