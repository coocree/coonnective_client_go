package api

import (
	"fmt"
	"time"
)

type ResponseModel struct {
	Success bool
	Errors  []ErrorModel
	Data    interface{}
}

func Response(success bool, errors []ErrorModel, data interface{}) ResponseModel {
	return ResponseModel{Success: success, Errors: errors, Data: data}
}

func (r ResponseModel) Result(endpoint string) interface{} {
	if r.Data != nil {
		response := r.Data.(map[string]interface{})
		if result, ok := response[endpoint]; ok {
			return result.(map[string]interface{})["result"]
		}
	}
	return nil
}

func (r ResponseModel) ResultError(endpoint string) interface{} {
	if r.Data != nil {
		response := r.Data.(map[string]interface{})
		if result, ok := response[endpoint]; ok {
			return result.(map[string]interface{})["error"]
		}
	}
	return nil
}

func (r ResponseModel) PageInfo(endpoint string) interface{} {
	if r.Data != nil {
		response := r.Data.(map[string]interface{})
		if result, ok := response[endpoint]; ok {
			return result.(map[string]interface{})["pageInfo"]
		}
	}
	return nil
}

func (r ResponseModel) Endpoint(name string) EndpointModel {
	if r.Data != nil {
		response := r.Data.(map[string]interface{})
		if endpoint, ok := response[name]; ok {
			return Endpoint(endpoint)
		}
	}

	apiError := ApiError([]string{"O endpoint '" + name + "' não foi encontrado no resultado da query/mutation de interação com a base de dados."}, "Response", "Response", "ENDPOINT_NOT_FOUND", time.Now(), r.Data)
	E(apiError.ToString(), apiError.Code, nil)
	//panic(apiError.Code)

	//TODO: throw exception
	fmt.Println(apiError.ToString())
	return EndpointModel{}
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
