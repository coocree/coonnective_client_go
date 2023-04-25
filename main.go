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

func mainXXXX() {

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

type DooSimucCreateResponse struct {
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

func DooSimucCreate() api.ResponseModel {
	variables := map[string]interface{}{"input": map[string]interface{}{
		"idSimuc": 1111444, "idSimcon": 2040001, "latitude": -23.59495, "longitude": -46.82484,
		"applicationType": "lightFixture", "numberOfLightFixtures": 1, "lightSensor": true, "dimmer": true,
		"contactCharacteristic": "nf", "sector": 2, "subSector": 3,
		"areaGroup": "LAPA", "observations": "Texto de observação", "label": "LAPA",
	}}
	graphQL := `
mutation DooSimucCreate($input: DooSimucCreateInput!) {
  DooSimucCreate(input: $input) {
    result {
      idSimuc
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

	return api.Dao(graphQL, variables, &DooSimucCreateResponse{})
}

func main() {

	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiaHR0cDovLzE5Mi4xNjguMS4zNDo0NjAwIl0sImNudCI6InN1cG9ydGUtYXBpQGFlZ2lzLmFwcC5iciIsImNpZCI6IjYyYWU1YmI5YmI1NTI5MDAyZGJmOTliMSIsImNuYSI6IkNvcmRhIEx1ei1HZW9Qb2ludCIsImN0eCI6ImNsaWVudCIsImV4cCI6NDc2NTk5MzkxMywiaWF0IjoxNjU1NTkzOTEzLCJpc3MiOiJodHRwOi8vMTkyLjE2OC4xLjM0OjQ2MDAiLCJzY29wZSI6Im11dDphdXRob3JpemUiLCJzdWIiOiJhdXRoX3Rva2VuIiwidWlkIjoic2Q1YTI4ZmM4ZCIsInByaiI6ImJhcnJhLWRvLWNvcmRhIn0.JMXFPI_JD9OTe_8r3Hk-orL3fRJjsvWZkUpCTa7yZQA"
	api.Connect(&token)

	/*apiResponse := kdlSimucEventReset()
	if !apiResponse.IsValid() {
		apiResponse.ThrowException()
	}
	apiEndpoint := apiResponse.Endpoint("KdlSimucEventReset")
	if !apiEndpoint.IsValid() {
		apiEndpoint.ThrowException()
	}

	fmt.Println(apiEndpoint.Result)*/

	apiResponse := DooSimucCreate()

	if !apiResponse.IsValid() {
		apiResponse.ThrowException()
	}

	apiEndpoint := apiResponse.Endpoint().(*DooSimucCreateResponse)
	dooSimucCreate := apiEndpoint.Data.DooSimucCreate

	if dooSimucCreate.Success == false {
		fmt.Println("dooSimucCreate.Error", dooSimucCreate.Error)
		return
	}

	fmt.Printf("%v: -->> %+v\n", "xxx--->>", dooSimucCreate.Result)
	fmt.Println("apiEndpoint.ResultXXXX", dooSimucCreate.Result.IdSimuc)

	//fmt.Println("apiEndpoint.Result", apiEndpoint.Result)

	//fmt.Println(apiEndpoint.Result.(map[string]interface{})["idSimuc"].(int32))

	//fmt.Printf("%v: -->> %+v\n", "apiEndpoint--->>", apiEndpoint)
}
