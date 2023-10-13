package coonective

import "github.com/coocree/coonnective_client_go/api"

type KdlHasConnectionFase1Response struct {
	Data struct {
		KdlHasConnectionFase1 struct {
			Result struct {
				IdLine int `bson:"idLine,omitempty" json:"idLine"`
			} `bson:"result,omitempty" json:"result"`
			Error       api.ErrorModel
			ElapsedTime string
			Success     bool
		}
	}
}

func KdlHasConnectionFase1() api.ResponseModel {
	variables := map[string]interface{}{
		"filter": map[string]interface{}{
			"codUser":  1,
			"idSimcon": 2206203,
		},
	}
	graphQL := `
query kdlHasConnectionFase1($filter: KdlHasConnectionFase1Filter!) {
	kdlHasConnectionFase1(filter: $filter) {
		result {
			idLine
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
