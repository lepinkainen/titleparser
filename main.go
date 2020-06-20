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

		decoder := json.NewDecoder(r.Body)
		var query lambda.TitleQuery
		err := decoder.Decode(&query)
		if err != nil {
			log.Errorln("OMGLOL")
			return
		}

		res, err := lambda.HandleRequest(context.Background(), query)
		if err != nil {
			_ = fmt.Errorf("Error handling request: %#v", err)
		}

		q, err := json.Marshal(res)
		if err != nil {
			_ = fmt.Errorf("Error marshaling response JSON: %#v", err)
		}

		fmt.Fprint(w, string(q))
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8081", nil))

}
