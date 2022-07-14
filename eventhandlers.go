package main

import (
	"log"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

// HandleTestTriggeredEvent handles test.triggered events
// TODO: add in your handler code
func HandleTestTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.TestTriggeredEventData) error {
	log.Printf("Handling test.triggered Event: %s", incomingEvent.Context.GetID())

	// check if action is supported
	if data.Test.TestStrategy == "functional" {
		// -----------------------------------------------------
		// 1. Send Action.Started Cloud-Event
		// -----------------------------------------------------
		myKeptn.SendTaskStartedEvent(data, ServiceName)

		// -----------------------------------------------------
		// 2. Implement your remediation action here
		// -----------------------------------------------------
		time.Sleep(5 * time.Second) // Example: Wait 5 seconds. Maybe the problem fixes itself.

		// -----------------------------------------------------
		// 3. Send Action.Finished Cloud-Event
		// -----------------------------------------------------
		myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
			Status:  keptnv2.StatusSucceeded, // alternative: keptnv2.StatusErrored
			Result:  keptnv2.ResultPass,      // alternative: keptnv2.ResultFailed
			Message: "Successfully sleeped!",
		}, ServiceName)

	} else {
		log.Printf("Retrieved unknown test strategy %s, skipping...", data.Test.TestStrategy)
		log.Printf("jainam-log the service name is %s", ServiceName)
		return nil
	}

	return nil
}
