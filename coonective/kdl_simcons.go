package coonective

import (
	"fmt"
	"github.com/coocree/coonnective_client_go/api"
)

type KdlSimconsAndTotalResponse struct {
	Data struct {
		KdlSimcons struct {
			Result []struct {
				IdSimcon  int     `bson:"idSimcon" json:"idSimcon"`
				Latitude  float64 `bson:"latitude" json:"latitude"`
				Longitude float64 `bson:"longitude" json:"longitude"`
				Icon      int     `bson:"icon" json:"icon"`
			} `bson:"result" json:"result"`
			Error       api.ErrorModel
			ElapsedTime string
			Success     bool
		}
		KdlSimconsTotal struct {
			Result []struct {
				IdSimcon int     `bson:"idSimcon" json:"idSimcon"`
				Total    float64 `bson:"total" json:"total"`
			} `bson:"result" json:"result"`
			Error       api.ErrorModel
			ElapsedTime string
			Success     bool
		}
	}
}

func KdlSimcons() api.ResponseModel {
	variables := map[string]interface{}{
		"kdlSimconsFilter":      map[string]interface{}{"codUser": 3},
		"kdlSimconsTotalFilter": map[string]interface{}{"codUser": 3},
	}
	graphQL := `
query kdlSimcons($kdlSimconsFilter: KdlSimconsFilter!, $kdlSimconsTotalFilter: KdlSimconsTotalFilter!) {
  kdlSimcons(filter: $kdlSimconsFilter) {
    result {
      icon
      idSimcon
      latitude
      longitude
    }
    success
    elapsedTime
    error {
      code
      createdAt
      messages
      module
      path
      variables
    }
  }
  kdlSimconsTotal(filter: $kdlSimconsTotalFilter) {
    result {
      idSimcon
      total
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

func KdlSimconsExec() {
	// Envie a consulta para a API GraphQL
	apiResponse := KdlSimcons()

	var response KdlSimconsAndTotalResponse
	apiResponse.Endpoint(&response)

	fmt.Println("KdlSimcons--->>>", response.Data.KdlSimcons.Result)
	fmt.Println(" ")
	fmt.Println(" ")
	fmt.Println("KdlSimconsTotal--->>>", response.Data.KdlSimconsTotal.Result)

	kdlSimcons := apiResponse.EndpointToMap("kdlSimconsTotal")
	fmt.Println("KdlSimcons--->>>", kdlSimcons)

	// Resultado da consulta
	/*apiEndpoint := apiResponse.Endpoint().(*KdlSimconsAndTotalResponse)
	kdlSimcons := apiEndpoint.Data.KdlSimcons
	fmt.Println("KdlSimcons--->>>", kdlSimcons.Result)

	kdlSimconsTotal := apiEndpoint.Data.KdlSimconsTotal
	fmt.Println("kdlSimconsTotal--->>>", kdlSimconsTotal.Result)*/
}
