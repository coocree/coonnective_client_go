package main

import (
	"coonnective_client_go/api"
	"fmt"
)

func EventReset() api.ResponseModel {
	variables := map[string]interface{}{}
	graphQL := `
mutation EventReset {
	EventReset(filter: {idEvent: "123"}) {
		result {
			idEvent
			status
		}
		error {
			code
			createdAt
			messages
			module
			path
			variables
		}
		elapsedTime
		success
	}
}
`
	return api.Dao(graphQL, variables)
}

func main() {

	token := ""
	api.Connect(&token)

	apiResponse := EventReset()
	if !apiResponse.IsValid() {
		apiResponse.ThrowException()
	}
	apiEndpoint := apiResponse.Endpoint("EventReset")
	if !apiEndpoint.IsValid() {
		apiEndpoint.ThrowException()
	}

	fmt.Println(apiEndpoint.Result)
}
