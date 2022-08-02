package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
	"github.com/kelseyhightower/envconfig"
	keptn "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

var keptnOptions = keptn.KeptnOpts{}

type gracefulShutdownKeyType struct{}

var gracefulShutdownKey = gracefulShutdownKeyType{}

// Opaque key type used for graceful shutdown context value
type keptnQuitType struct{}

var serviceRunnerQuit = keptnQuitType{}

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int `envconfig:"RCV_PORT" default:"8080"`
	// Path to which cloudevents are sent
	Path string `envconfig:"RCV_PATH" default:"/"`
	// Whether we are running locally (e.g., for testing) or on production
	Env string `envconfig:"ENV" default:"local"`
	// URL of the Keptn configuration service (this is where we can fetch files from the config repo)
	ConfigurationServiceUrl string `envconfig:"CONFIGURATION_SERVICE" default:""`
}

// ServiceName specifies the current services name (e.g., used as source when sending CloudEvents)
const ServiceName = "k6-service"

/**
 * Parses a Keptn Cloud Event payload (data attribute)
 */
func parseKeptnCloudEventPayload(event cloudevents.Event, data interface{}) error {
	err := event.DataAs(data)
	if err != nil {
		log.Fatalf("Got Data Error: %s", err.Error())
		return err
	}
	return nil
}

/**
 * This method gets called when a new event is received from the Keptn Event Distributor
 * Depending on the Event Type will call the specific event handler functions, e.g: handleDeploymentFinishedEvent
 * See https://github.com/keptn/spec/blob/0.2.0-alpha/cloudevents.md for details on the payload
 */
func processKeptnCloudEvent(ctx context.Context, event cloudevents.Event) error {
	// create keptn handler
	log.Printf("Initializing Keptn Handler")
	myKeptn, err := keptnv2.NewKeptn(&event, keptnOptions)
	if err != nil {
		return errors.New("Could not create Keptn Handler: " + err.Error())
	}

	log.Printf("gotEvent(%s): %s - %s", event.Type(), myKeptn.KeptnContext, event.Context.GetID())

	if err != nil {
		log.Printf("failed to parse incoming cloudevent: %v", err)
		return err
	}

	switch event.Type() {
	// -------------------------------------------------------
	// sh.keptn.event.test
	case keptnv2.GetTriggeredEventType(keptnv2.TestTaskName): // sh.keptn.event.test.triggered
		log.Printf("Processing Test.Triggered Event")

		eventData := &keptnv2.TestTriggeredEventData{}
		parseKeptnCloudEventPayload(event, eventData)

		return HandleTestTriggeredEvent(myKeptn, event, eventData)

	}

	// Unknown Event -> Throw Error!
	var errorMsg string = fmt.Sprintf("Unhandled Keptn Cloud Event: %s", event.Type())

	log.Print(errorMsg)
	return errors.New(errorMsg)
}

/**
 * Usage: ./main
 * no args: starts listening for cloudnative events on localhost:port/path
 *
 * Environment Variables
 * env=runlocal   -> will fetch resources from local drive instead of configuration service
 */
func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}

	os.Exit(_main(os.Args[1:], env))
}

/**
 * Opens up a listener on localhost:port/path and passes incoming requets to gotEvent
 */
func _main(args []string, env envConfig) int {
	// configure keptn options
	if env.Env == "local" {
		log.Println("env=local: Running with local filesystem to fetch resources")
		keptnOptions.UseLocalFileSystem = true
	}

	keptnOptions.ConfigurationServiceURL = env.ConfigurationServiceUrl

	log.Println("Starting k6-service...")
	log.Printf("    on Port = %d; Path=%s", env.Port, env.Path)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	log.Printf("Creating new http handler")

	// configure http server to receive cloudevents
	p, err := cloudevents.NewHTTP(
		cloudevents.WithPath(env.Path), cloudevents.WithPort(env.Port), cloudevents.WithGetHandlerFunc(HTTPGetHandler),
	)

	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	err = c.StartReceiver(ctx, processKeptnCloudEvent)
	if err != nil {
		log.Fatalf("CloudEvent receiver stopped with error: %v", err)
	}
	log.Printf("Shutdown complete.")
	return 0
}

// HTTPGetHandler will handle all requests for '/health' and '/ready'
func HTTPGetHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/health":
		healthEndpointHandler(w, r)
	case "/ready":
		healthEndpointHandler(w, r)
	default:
		endpointNotFoundHandler(w, r)
	}
}

// HealthHandler reports a basic health check back
func healthEndpointHandler(w http.ResponseWriter, r *http.Request) {
	type StatusBody struct {
		Status string `json:"status"`
	}

	status := StatusBody{Status: "OK"}

	body, _ := json.Marshal(status)

	w.Header().Set("content-type", "application/json")

	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

// endpointNotFoundHandler will return 404 for requests
func endpointNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	type StatusBody struct {
		Status string `json:"status"`
	}

	status := StatusBody{Status: "NOT FOUND"}

	body, _ := json.Marshal(status)

	w.Header().Set("content-type", "application/json")

	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}
