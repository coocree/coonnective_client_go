package api

import (
	"fmt"
	"log"
)

func E(message interface{}, err string, stackTrace *string) {
	logMessage := fmt.Sprintf("%v\nError: %v", message, err)
	if stackTrace != nil {
		logMessage = fmt.Sprintf("%v\nStack Trace: %v", logMessage, stackTrace)
	}
	log.Println("ERROR: ", logMessage)
}
