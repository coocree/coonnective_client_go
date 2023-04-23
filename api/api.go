package api

import "strings"

func Dao(graphQL string, variables map[string]interface{}) ResponseModel {
	apiParams := Params(graphQL)
	if errParams := apiParams.IsValid(); errParams != nil {
		return ResponseModel{
			Errors: []ErrorModel{*errParams},
		}
	}

	apiConnect, errConnect := Connect(nil)
	if errConnect != nil {
		return ResponseModel{
			Errors: []ErrorModel{*errConnect},
		}
	}

	if strings.ToLower(apiParams.Module) == "query" {
		return apiConnect.Query(graphQL, variables)
	}
	return apiConnect.Mutation(graphQL, variables)
}
