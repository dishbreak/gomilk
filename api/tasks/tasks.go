package tasks

import (
	"encoding/json"
	"time"

	"github.com/dishbreak/gomilk/api"
	log "github.com/sirupsen/logrus"
)

type tagsRecord struct {
	Tag []string
}

type tagsWrapper struct {
	tagsRecord
	Empty bool
}

// UnmarshalJSON handles the situation where tasks is an empty array when empty.
func (t *tagsWrapper) UnmarshalJSON(data []byte) error {
	// GRRR. When tags is empty, it's a JSON array.
	if string(data) == "[]" {
		t.Empty = true
		return nil
	}

	// And when it's not, it's an object.
	return json.Unmarshal(data, &t.tagsRecord)
}

// TaskResponse covers what RTM calls a "taskseries"
type TaskResponse struct {
	RawID    string `json:"id"`
	Created  time.Time
	Modified time.Time
	RawName  string `json:"name"`
	Source   string
	URL      string
	Tags     tagsWrapper
	Task     []struct {
		RawID      string `json:"id"`
		Due        string
		HasDueTime string `json:"has_due_time"`
		Added      time.Time
		Completed  string
		Deleted    string
		Priority   string
		Postponed  string
		Estimate   string
	}
}

// TaskList encapsulates a List API object.
type TaskList struct {
	ID         string
	Taskseries []TaskResponse
}

// TaskAddResponse contains the response for the call to rtm.tasks.add
type TaskAddResponse struct {
	Rsp struct {
		Stat        string
		Tranasction api.TransactionRecord
		List        TaskList
	}
}

// TaskCompleteResponse contains the response for the call to rtm.tasks.complete
type TaskCompleteResponse TaskAddResponse

// TaskSetDueDateResponse contains the response for the call to
// rtm.tasks.setDueDate
type TaskSetDueDateResponse TaskAddResponse

// GetListResponse contains the response for the call to rtm.tasks.getList
type GetListResponse struct {
	Rsp struct {
		Stat  string
		Tasks struct {
			List []TaskList
		}
	}
}

/*
Add invokes rtm.task.add to create a new task. This method uses Smart Add and can only add top-level tasks.
*/
func Add(apiToken string, name string, timelineID string) (TaskAddResponse, error) {
	args := map[string]string{
		"api_key":    api.APIKey,
		"auth_token": apiToken,
		"timeline":   timelineID,
		"parse":      "1", // parse name using Smart Add
		"name":       name,
	}

	var response TaskAddResponse
	unmarshal := func(b []byte) error {
		return json.Unmarshal(b, &response)
	}

	err := api.GetMethod("rtm.tasks.add", args, unmarshal)
	if err != nil {
		return response, err
	}

	return response, nil

}

// GetList invokes rtm.tasks.getList to return a listing of lists, taskseries, and tasks.
func GetList(apiToken string, filter string) (GetListResponse, error) {
	args := map[string]string{
		"api_key":    api.APIKey,
		"auth_token": apiToken,
	}

	if filter != "" {
		args["filter"] = filter
	}

	var response GetListResponse
	unmarshal := func(b []byte) error {
		return json.Unmarshal(b, &response)
	}

	err := api.GetMethod("rtm.tasks.getList", args, unmarshal)
	if err != nil {
		return response, err
	}

	log.WithFields(log.Fields{
		"response": response,
	}).Debug("")

	return response, nil

}

// Complete invokes rtm.tasks.complete
func Complete(apiToken, timelineID, listID, taskseriesID, taskID string) (TaskCompleteResponse, error) {
	args := map[string]string{
		"api_key":       api.APIKey,
		"auth_token":    apiToken,
		"timeline":      timelineID,
		"list_id":       listID,
		"taskseries_id": taskseriesID,
		"task_id":       taskID,
	}

	var response TaskCompleteResponse
	unmarshal := func(b []byte) error {
		return json.Unmarshal(b, &response)
	}

	err := api.GetMethod("rtm.tasks.complete", args, unmarshal)
	if err != nil {
		return response, err
	}

	return response, nil
}

// SetDueDate invokes rtm.tasks.setDueDate
func SetDueDate(apiToken, timelineID, listID, taskseriesID, taskID, due string) (TaskSetDueDateResponse, error) {
	args := map[string]string{
		"api_key":       api.APIKey,
		"auth_token":    apiToken,
		"timeline":      timelineID,
		"list_id":       listID,
		"taskseries_id": taskseriesID,
		"task_id":       taskID,
		"due":           due,
		"parse":         "1", // always parse the due date
	}

	var response TaskSetDueDateResponse
	unmarshal := func(b []byte) error {
		return json.Unmarshal(b, &response)
	}

	err := api.GetMethod("rtm.tasks.setDueDate", args, unmarshal)
	if err != nil {
		return response, err
	}

	return response, nil
}
