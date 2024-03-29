package main

import (
	"encoding/json"
	"fmt"
	"github.com/coocree/coonnective_client_go/api"
	"github.com/coocree/coonnective_client_go/coonective"
	"log"
	"nhooyr.io/websocket"
	"time"
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
	//return api.Dao(graphQL, variables, &EventResetResponse{})
	return api.Dao(graphQL, variables)
}

func mainXXXX() {

	token := ""
	api.Connect(&token)

	apiResponse := EventReset()
	//apiEndpoint := apiResponse.Endpoint().(*EventResetResponse)
	var response EventResetResponse
	apiResponse.Endpoint(&response)
	/*dooSimucCreate := apiEndpoint.Data.DooSimucCreate

	if !dooSimucCreate.Success {
		fmt.Println(dooSimucCreate.Error.Messages)
	}

	fmt.Println(dooSimucCreate.Result)*/
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

	//return api.Dao(graphQL, variables, &DooSimucCreateResponse{})
	return api.Dao(graphQL, variables)
}

type PostMessageRespose struct {
	Data struct {
		PostMessage string
	}
}

func PostMessage() api.ResponseModel {
	variables := map[string]interface{}{"input": map[string]interface{}{
		"idSimuc": 1111444, "idSimcon": 2040001, "latitude": -23.59495, "longitude": -46.82484,
		"applicationType": "lightFixture", "numberOfLightFixtures": 1, "lightSensor": true, "dimmer": true,
		"contactCharacteristic": "nf", "sector": 2, "subSector": 3,
		"areaGroup": "LAPA", "observations": "Texto de observação", "label": "LAPA",
	}}
	graphQL := `
mutation MyMutation {
  postMessage(content: "aaaa", user: "123")
}
`

	//return api.Dao(graphQL, variables, &DooSimucCreateResponse{})
	return api.Dao(graphQL, variables)
}

type KdlSimconHasConnectionResponse struct {
	Data struct {
		KdlSimconHasConnection struct {
			Result struct {
				Status string
			}
			Error       api.ErrorModel
			ElapsedTime string
			Success     bool
		}
	}
}

func KdlSimconHasConnection() api.ResponseModel {
	variables := map[string]interface{}{"filter": map[string]interface{}{"idSimcon": 2040001}}
	graphQL := `
query KdlSimconHasConnection($filter: KdlSimconHasConnectionFilter!) {
  KdlSimconHasConnection(filter: $filter) {
    result {
      status
    }
    error {
      variables
      path
      module
      messages
      createdAt
      code
    }
    elapsedTime
    success
  }
}
`
	//return api.Dao(graphQL, variables, &KdlSimconHasConnectionResponse{})
	return api.Dao(graphQL, variables)
}

type KdlCustomerResponse struct {
	Data struct {
		KdlCustomer struct {
			Result struct {
				Address      string
				City         string
				CodUsr       int32
				Country      string
				CreatedAt    time.Time
				Users        int32
				User         string
				TelephoneTwo string
				TelephoneOne string
				State        string
				Simucs       int32
				Simcons      int32
				PostalCode   string
				Passwd       string
				Neighborhood string
				Email        string
			}
			Error       api.ErrorModel
			ElapsedTime string
			Success     bool
		}
	}
}

func KdlCustomer() api.ResponseModel {
	variables := map[string]interface{}{"filter": map[string]interface{}{"codusr": 10}}
	graphQL := `
query kdlCustomer($filter: KdlCustomerFilter!) {
  kdlCustomer(filter: $filter) {
    error {
      code
      createdAt
      messages
      module
      path
      variables
    }
    elapsedTime
    result {
      address
      city
      codUsr
      country
      createdAt
      users
      user
      telephoneTwo
      telephoneOne
      state
      simucs
      simcons
      postalCode
      passwd
      neighborhood
      email
    }
    success
  }
}
`
	//return api.Dao(graphQL, variables, &KdlCustomerResponse{})
	return api.Dao(graphQL, variables)
}

func mainDDD() {

	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiaHR0cDovLzE5Mi4xNjguMS4zNDo0NjAwIl0sImNudCI6InN1cG9ydGUtYXBpQGFlZ2lzLmFwcC5iciIsImNpZCI6IjYyYWU1YmI5YmI1NTI5MDAyZGJmOTliMSIsImNuYSI6IkNvcmRhIEx1ei1HZW9Qb2ludCIsImN0eCI6ImNsaWVudCIsImV4cCI6NDc2NTk5MzkxMywiaWF0IjoxNjU1NTkzOTEzLCJpc3MiOiJodHRwOi8vMTkyLjE2OC4xLjM0OjQ2MDAiLCJzY29wZSI6Im11dDphdXRob3JpemUiLCJzdWIiOiJhdXRoX3Rva2VuIiwidWlkIjoic2Q1YTI4ZmM4ZCIsInByaiI6ImJhcnJhLWRvLWNvcmRhIn0.JMXFPI_JD9OTe_8r3Hk-orL3fRJjsvWZkUpCTa7yZQA"
	api.Connect(&token)

	apiResponse := KdlCustomer()

	//apiEndpoint := apiResponse.Endpoint().(*KdlCustomerResponse)
	var response KdlCustomerResponse
	apiResponse.Endpoint(&response)
	/*kdlCustomer := apiEndpoint.Data.KdlCustomer

	if kdlCustomer.Success == false {
		fmt.Println("hasConnection.Error", kdlCustomer.Error)
		return
	}

	fmt.Println("Result", kdlCustomer.Result)

	fmt.Println("ElapsedTime", kdlCustomer.ElapsedTime)
	fmt.Println("Success", kdlCustomer.Success)
	fmt.Println("Error", kdlCustomer.Error)*/

	/*apiResponse := KdlSimconHasConnection()
	if !apiResponse.IsValid() {
		apiResponse.ThrowException()
	}

	apiEndpoint := apiResponse.Endpoint().(*KdlSimconHasConnectionResponse)
	hasConnection := apiEndpoint.Data.KdlSimconHasConnection

	if hasConnection.Success == false {
		fmt.Println("hasConnection.Error", hasConnection.Error)
		return
	}

	fmt.Println("Result", hasConnection.Result.Status)

	fmt.Println("ElapsedTime", hasConnection.ElapsedTime)
	fmt.Println("Success", hasConnection.Success)
	fmt.Println("Error", hasConnection.Error)*/

	/*apiResponse := DooSimucCreate()
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
	fmt.Println("apiEndpoint.ResultXXXX", dooSimucCreate.Result.IdSimuc)*/

}

type Response struct {
	Data struct {
		DOOMessage struct {
			Success     bool   `json:"success"`
			ElapsedTime string `json:"elapsedTime"`
			Result      struct {
				ID         string `json:"id"`
				User       string `json:"user"`
				Content    string `json:"content"`
				InsertedID string `json:"insertedId"`
			} `json:"result"`
			Error interface{} `json:"error"`
		} `json:"dooMessage"`
	} `json:"data"`
}

func resultSubscription(rawMsg json.RawMessage) {
	// Decodificar a parte específica do JSON usando RawMessage
	var response Response
	if err := json.Unmarshal(rawMsg, &response); err != nil {
		fmt.Println("Erro ao decodificar Person:", err)
		return
	}

	fmt.Println("Nome:", response.Data)
	fmt.Println("Idade:", response.Data.DOOMessage.Result.ID)

}

func mainWSA() {
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	//--------------------------------------------------------

	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiaHR0cDovLzE5Mi4xNjguMS4zNDo0NjAwIl0sImNudCI6InN1cG9ydGUtYXBpQGFlZ2lzLmFwcC5iciIsImNpZCI6IjYyYWU1YmI5YmI1NTI5MDAyZGJmOTliMSIsImNuYSI6IkNvcmRhIEx1ei1HZW9Qb2ludCIsImN0eCI6ImNsaWVudCIsImV4cCI6NDc2NTk5MzkxMywiaWF0IjoxNjU1NTkzOTEzLCJpc3MiOiJodHRwOi8vMTkyLjE2OC4xLjM0OjQ2MDAiLCJzY29wZSI6Im11dDphdXRob3JpemUiLCJzdWIiOiJhdXRoX3Rva2VuIiwidWlkIjoic2Q1YTI4ZmM4ZCIsInByaiI6ImJhcnJhLWRvLWNvcmRhIn0.JMXFPI_JD9OTe_8r3Hk-orL3fRJjsvWZkUpCTa7yZQA"
	api.Connect(&token)

	apiResponse := coonective.WkfEvents()

	//apiEndpoint := apiResponse.Endpoint().(*coonective.WkfEventsResponse)
	var response coonective.WkfEventsResponse
	apiResponse.Endpoint(&response)
	/*wkfEvents := apiEndpoint.Data.WkfEvents

	if wkfEvents.Success == false {
		fmt.Println("hasConnection.Error", wkfEvents.Error)
		return
	}

	fmt.Println("Result--->>>", wkfEvents.Result)

	for _, result := range wkfEvents.Result {
		fmt.Println("Result--->>>", result.Workers)
		for _, worker := range result.Workers {
			fmt.Println("Result--->>>", worker)
		}
	}

	fmt.Println("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")*/
}

func mainXXX() {
	url := "ws://localhost:4600/coonective"
	fmt.Println("url: ", url)

	conn, err := api.ConnectToGraphQLServer(url)
	if err != nil {
		log.Fatal("Error connecting to GraphQL server:", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	subscription := `
		subscription dooMessage {
    dooMessage {
        success
        elapsedTime
        result {
            id
            user
            content
            insertedId
        }
        error {
            createdAt
            messages
            module
            path
            code
            variables
        }
    }
}
	`

	err = api.Subscribe(conn, subscription)
	if err != nil {
		log.Fatal("Error subscribing:", err)
	}

	api.ReceiveMessages(conn, resultSubscription)
}

func main() {
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiaHR0cDovLzE5Mi4xNjguMS4zNDo0NjAwIl0sImNudCI6InN1cG9ydGUtYXBpQGFlZ2lzLmFwcC5iciIsImNpZCI6IjYyYWU1YmI5YmI1NTI5MDAyZGJmOTliMSIsImNuYSI6IkNvcmRhIEx1ei1HZW9Qb2ludCIsImN0eCI6ImNsaWVudCIsImV4cCI6NDc2NTk5MzkxMywiaWF0IjoxNjU1NTkzOTEzLCJpc3MiOiJodHRwOi8vMTkyLjE2OC4xLjM0OjQ2MDAiLCJzY29wZSI6Im11dDphdXRob3JpemUiLCJzdWIiOiJhdXRoX3Rva2VuIiwidWlkIjoic2Q1YTI4ZmM4ZCIsInByaiI6ImJhcnJhLWRvLWNvcmRhIn0.JMXFPI_JD9OTe_8r3Hk-orL3fRJjsvWZkUpCTa7yZQA"
	api.Connect(&token)

	fmt.Println("xxxxx")
	// Envie a consulta para a API GraphQL
	coonective.KdlSimconsExec()

	/*if !apiResponse.IsValid() {
		apiResponse.ThrowException()
	}

	// Resultado da consulta
	apiEndpoint := apiResponse.Endpoint().(*coonective.KdlHasConnectionFase1Response)
	wkfEvents := apiEndpoint.Data.KdlHasConnectionFase1
	fmt.Println(wkfEvents.Result.IdLine)*/
}
