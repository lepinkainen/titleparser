package handler

import "github.com/lepinkainen/titleparser/lambda"

// YleAreena handler TBD
func YleAreena(url string) (string, error) {
	return "Areena custom handler", nil
}

func init() {
	lambda.RegisterHandler(".*?areena.yle.fi/.*", YleAreena)
}
