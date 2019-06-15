package tasks

import (
	"encoding/json"
	"time"

	"github.com/dishbreak/gomilk/api"
)

// TaskAddResponse contains the response for the call to rtm.tasks.add
type TaskAddResponse struct {
	Rsp struct {
		Stat        string
		Transaction struct {
			ID       string `json:"id"`
			Undoable string
		}
		List struct {
			ID         string `json:"id"`
			Taskseries []struct {
				ID       string `json:"id"`
				Created  time.Time
				Modified time.Time
				Name     string
				Source   string
				Task     []struct {
					ID  string `json:"id"`
					Due string
				}
			}
		}
	}
}

type Task interface {
	Name() string
	Due() *time.Time
}

/*
Name returns the task name.
*/
func (t *TaskAddResponse) Name() string {
	return t.Rsp.List.Taskseries[0].Name
}

/*
Due returns the due date.
*/
func (t *TaskAddResponse) Due() *time.Time {
	parsed, err := time.Parse(time.RFC3339, t.Rsp.List.Taskseries[0].Task[0].Due)
	if err != nil {
		return nil
	}
	return &parsed
}

type shellTask struct {
	name string
}

/*
Add invokes rtm.task.add to create a new task. This method uses Smart Add and can only add top-level tasks.
*/
func Add(apiToken string, name string, timelineID string) (Task, error) {
	args := map[string]string{
		"api_key":    api.APIKey,
		"auth_token": apiToken,
		"timeline":   timelineID,
		"parse":      "1", // parse name using Smart Add
		"name":       name,
	}

	var result Task
	var response TaskAddResponse
	unmarshal := func(b []byte) error {
		return json.Unmarshal(b, &response)
	}

	err := api.GetMethod("rtm.tasks.add", args, unmarshal)
	if err != nil {
		return nil, err
	}

	if result == nil {
		result = &response
	}

	return result, nil

}
