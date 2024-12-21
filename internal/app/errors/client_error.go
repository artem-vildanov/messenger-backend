package errors

import "log"

func ClientConnectionPanic(client, msg string) {
	log.Fatalf("Failed to establish connection with %s: %s", client, msg)
}