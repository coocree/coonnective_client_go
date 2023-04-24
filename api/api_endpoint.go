package api

import (
	"fmt"
	"time"
)

type EndpointModel struct {
	Success     bool
	Error       *ErrorModel
	Result      interface{}
	ElapsedTime string
}

func Endpoint(value interface{}) EndpointModel {
	if valueMap, ok := value.(map[string]interface{}); ok {
		if valueMap["success"] == nil {
			apiError := ApiError([]string{"O parametro 'success' não foi encontrado no resultado da query/mutation de interação com a base de dados."},
				"Endpoint", "Endpoint", "SUCCESS_NOT_FOUND", time.Now(), value)
			E(apiError.ToString(), apiError.Code, nil)
			panic(apiError.Code)
		}

		if valueMap["error"] == nil && valueMap["result"] == nil {
			apiError := ApiError([]string{"O parametro 'result' não foi encontrado no resultado da query/mutation de interação com a base de dados."},
				"Endpoint", "Endpoint", "RESULT_NOT_FOUND", time.Now(), value)
			E(apiError.ToString(), apiError.Code, nil)
			panic(apiError.Code)
		}

		endpoint := EndpointModel{
			Success:     valueMap["success"].(bool),
			Result:      valueMap["result"],
			ElapsedTime: valueMap["elapsedTime"].(string),
		}

		if valueMap["error"] != nil {
			errorMap := valueMap["error"].(map[string]interface{})
			endpoint.Error = &ErrorModel{
				CreatedAt: errorMap["createdAt"].(string),
				Code:      errorMap["code"].(string),
				Path:      errorMap["path"].(string),
				Messages:  []string{fmt.Sprintf("%v", errorMap["messages"])},
				Module:    errorMap["module"].(string),
				Variables: errorMap["variables"],
			}
		}

		return endpoint
	}

	return EndpointModel{}
}

func (e EndpointModel) IsValid() bool {
	return e.Success && e.Result != nil && e.Error == nil
}

func (e EndpointModel) ThrowException() {
	if e.Error != nil {
		E(e.Error.ToString(), e.Error.Code, nil)
		panic(e.Error.Code)
	}
}

func (e EndpointModel) ToString() string {
	return fmt.Sprintf("Instance of Endpoint(result:%v, success:%v, error:%v)", e.Result, e.Success, e.Error)
}
