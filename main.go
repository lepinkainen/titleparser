package main

import (
	// fake import for handlers to run their init() functions
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	_ "github.com/lepinkainen/titleparser/handler"
	"github.com/lepinkainen/titleparser/lambda"

	awslambda "github.com/aws/aws-lambda-go/lambda"
)

func main() {

	var runmode = os.Getenv("RUNMODE")
	if runmode != "local" {
		awslambda.Start(lambda.HandleRequest)
		os.Exit(0)
	}

	fmt.Println("Running in local mode")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	/*
		TODO: Register custom handlers
		TODO: Fall back to default if no custom handler matches url
		TODO: JSON Protocol for communication, body payload w/ apikey
	*/

	http.HandleFunc("/title", func(w http.ResponseWriter, r *http.Request) {

		query := lambda.TitleQuery{
			URL: "https://imgur.com/gallery/NsM4oor",
		}

		res, err := lambda.HandleRequest(context.Background(), query)
		if err != nil {
			_ = fmt.Errorf("Error handling request: %#v", err)
		}

		q, err := json.Marshal(res)
		if err != nil {
			_ = fmt.Errorf("Error marshaling response JSON: %#v", err)
		}

		fmt.Fprintf(w, string(q))
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8081", nil))

}
