package main

import (
	"fmt"
	"github.com/coocree/coonnective_client_go/api"
)

type EventResetResponse struct {
	Data struct {
		DooSimucCreate struct {
			Result struct {
				IdSimuc int32
			}
			Error       api.ErrorModel
			ElapsedTime string
			Success     bool
		}
	}
}

func EventReset() api.ResponseModel {
	variables := map[string]interface{}{
		"idEvent": "123",
	}
	graphQL := `
mutation EventReset(input: EventResetInput!) {
	EventReset(filter: $input) {
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
	return api.Dao(graphQL, variables, &EventResetResponse{})
}

func main() {

	token := ""
	api.Connect(&token)

	apiResponse := EventReset()
	if !apiResponse.IsValid() {
		apiResponse.ThrowException()
	}
	apiEndpoint := apiResponse.Endpoint().(*EventResetResponse)
	dooSimucCreate := apiEndpoint.Data.DooSimucCreate

	if !dooSimucCreate.Success {
		fmt.Println(dooSimucCreate.Error.Messages)
	}

	fmt.Println(dooSimucCreate.Result)
}
