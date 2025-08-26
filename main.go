package main

import (
	// fake import for handlers to run their init() functions
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	_ "github.com/lepinkainen/titleparser/handler"
	"github.com/lepinkainen/titleparser/lambda"

	awslambda "github.com/aws/aws-lambda-go/lambda"
)

func main() {

	// TODO: API-key for security in local mode
	// TODO: Make port configurable

	var runmode = os.Getenv("RUNMODE")
	if runmode != "local" && runmode != "stdin" {
		awslambda.Start(lambda.HandleRequest)
		os.Exit(0)
	}

	if runmode == "stdin" {
		fmt.Println("Running in stdin mode")

		decoder := json.NewDecoder(os.Stdin)
		var query lambda.TitleQuery
		err := decoder.Decode(&query)
		if err != nil {
			log.Errorf("Error decoding JSON from stdin: %v", err)
			os.Exit(1)
		}

		res, err := lambda.HandleRequest(context.Background(), query)
		if err != nil {
			log.Errorf("Error handling request: %v", err)
			os.Exit(1)
		}

		output, err := json.Marshal(res)
		if err != nil {
			log.Errorf("Error marshaling response JSON: %v", err)
			os.Exit(1)
		}

		fmt.Println(string(output))
		os.Exit(0)
	}

	fmt.Println("Running in local mode")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path)); err != nil {
			log.Warnf("Failed to write response: %v", err)
		}
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "Hi"); err != nil {
			log.Warnf("Failed to write response: %v", err)
		}
	})

	/*
		TODO: Register custom handlers
		TODO: Fall back to default if no custom handler matches url
		TODO: JSON Protocol for communication, body payload w/ apikey
	*/

	// Local mode calls like this with httpie:
	// http http://localhost:8081/title url="http://mantta.fi"

	http.HandleFunc("/title", func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var query lambda.TitleQuery
		err := decoder.Decode(&query)
		if err != nil {
			log.Errorln("No URL given")
			return
		}

		res, err := lambda.HandleRequest(context.Background(), query)
		if err != nil {
			_ = fmt.Errorf("error handling request: %#v", err)
		}

		q, err := json.Marshal(res)
		if err != nil {
			_ = fmt.Errorf("error marshaling response JSON: %#v", err)
		}

		if _, err := fmt.Fprint(w, string(q)); err != nil {
			log.Warnf("Failed to write response: %v", err)
		}
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8081", nil))

}
