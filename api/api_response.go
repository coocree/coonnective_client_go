package api

import (
	"encoding/json"
	"fmt"
)

type ResponseModel struct {
	Success bool
	Errors  []ErrorModel
	Data    interface{}
	Body    []byte
}

/*
token := apiResponse.EndpointMap("AclToken")

	if !token.IsValid() {
		token.ThrowException()
	}

	tokenResult := token.Result.(map[string]interface{})
	a.tokenType = tokenResult["tokenType"].(string)
*/

func (r ResponseModel) EndpointToMap(name any) map[string]interface{} {

	var data map[string]interface{}
	errJson := json.Unmarshal(r.Body, &data)
	if errJson != nil {
		fmt.Println("Mutation Error parsing JSON response: ", errJson)
	}

	if data != nil && name != nil {
		response := data["data"].(map[string]interface{})
		if endpoint, ok := response[name.(string)]; ok {
			return endpoint.(map[string]interface{})
		}
	}
	return data
}

func (r ResponseModel) Endpoint(response interface{}) {
	errJson := json.Unmarshal(r.Body, &response)
	if errJson != nil {
		fmt.Println("errJson: ", errJson)
	}
}

func (r ResponseModel) IsValid() bool {
	return r.Success && len(r.Errors) == 0
}

func (r ResponseModel) ThrowException() {
	if len(r.Errors) > 0 {
		for _, item := range r.Errors {
			E(item.ToString(), item.Code, nil)
			panic(item.Code)
		}
	}
}

func (r ResponseModel) ToString() string {
	return fmt.Sprintf("Instance of Response(data:%v, success:%v, errors:%v)", r.Data, r.Success, r.Errors)
}
