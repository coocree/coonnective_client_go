package coonective

import (
	"github.com/coocree/coonnective_client_go/api"
	"time"
)

type WkfEventActionResult struct {
	Index           int    `bson:"index,omitempty" json:"index"`
	Action          string `bson:"action,omitempty" json:"action"`
	Status          string `bson:"status,omitempty" json:"status"`
	NumInteractions int    `bson:"numInteractions,omitempty" json:"numInteractions"`
	Error           string `bson:"error,omitempty" json:"error"`
}

type WkfEventStepResult struct {
	Index      int                     `bson:"index,omitempty" json:"index"`
	Actions    []*WkfEventActionResult `bson:"actions,omitempty" json:"actions"`
	Status     string                  `bson:"status,omitempty" json:"status"`
	StepFailed *string                 `bson:"stepFailed,omitempty" json:"stepFailed"`
}

type WkfEventWorkerResult struct {
	Index  int                   `bson:"index,omitempty" json:"index"`
	Status string                `bson:"status,omitempty" json:"status"`
	Steps  []*WkfEventStepResult `bson:"steps,omitempty" json:"steps"`
}

type WkfEventsResponse struct {
	Data struct {
		WkfEvents struct {
			Result struct {
				ID              string                  `bson:"_id,omitempty" json:"_id"`
				Status          string                  `bson:"status,omitempty" json:"status"`
				CreatedAt       time.Time               `bson:"createdAt,omitempty" json:"createdAt"`
				ChangedAt       time.Time               `bson:"changedAt,omitempty" json:"changedAt"`
				NextAt          time.Time               `bson:"nextAt,omitempty" json:"nextAt"`
				Type            string                  `bson:"type,omitempty" json:"type"`
				NumInteractions int                     `bson:"numInteractions,omitempty" json:"numInteractions"`
				Workers         []*WkfEventWorkerResult `bson:"workers,omitempty" json:"workers"`
				EventFailed     *string                 `bson:"eventFailed,omitempty" json:"eventFailed"`
				Erros           []*string               `bson:"erros,omitempty" json:"erros"`
			}
			Error       api.ErrorModel
			ElapsedTime string
			Success     bool
		}
	}
}

func WkfEvents() api.ResponseModel {
	variables := map[string]interface{}{"filter": map[string]interface{}{}}
	graphQL := `
query wkfEvents {
  wkfEvents {
    result {
      _id
      nextAt
      eventFailed
      numInteractions
      status
      type
      workers {
        index
        status
        steps {
          index
          status
          stepFailed
          actions {
            action
            error
            index
            numInteractions
            status
          }
        }
      }
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
	return api.Dao(graphQL, variables, &WkfEventsResponse{})
}
