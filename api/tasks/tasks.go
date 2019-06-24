package tasks

import (
	"encoding/json"
	"time"

	"github.com/dishbreak/gomilk/api"
	"github.com/dishbreak/gomilk/model/task"
	log "github.com/sirupsen/logrus"
)

type taskResponse struct {
	RawID    string `json:"id"`
	Created  time.Time
	Modified time.Time
	RawName  string `json:"name"`
	Source   string
	URL      string
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

// TaskAddResponse contains the response for the call to rtm.tasks.add
type TaskAddResponse struct {
	Rsp struct {
		Stat        string
		Tranasction api.TransactionRecord
		List        struct {
			ID         string
			Taskseries []taskResponse
		}
	}
}

/*
Name returns the task name.
*/
func (t taskResponse) Name() string {
	return t.RawName
}

/*
DueDate returns the due date, or an error if it doesn't have one.
*/
func (t taskResponse) DueDate() (task.DateTime, error) {
	// Determine the timezone.
	zone, err := time.LoadLocation("Local")
	if err != nil {
		zone, err = time.LoadLocation("UTC")
		if err != nil {
			panic(err)
		}
	}

	parsed, err := time.Parse(time.RFC3339, t.Task[0].Due)
	if err != nil {
		parsed = time.Now().In(zone)
	}
	return task.DateTime{parsed.In(zone), t.DueDateHasTime()}, err
}

/*
DueDateHasTime returns true when the user gave a specific date and time for the due date,
false when the task just has a date.
*/
func (t taskResponse) DueDateHasTime() bool {
	return t.Task[0].HasDueTime == "1"
}

/*
IsCompleted returns true when the completed field is a datetime.
*/
func (t taskResponse) IsCompleted() bool {
	completed := true
	_, err := time.Parse(time.RFC3339, t.Task[0].Completed)
	if err != nil {
		completed = false
	}

	return completed
}

func (t taskResponse) ID() string {
	return t.RawID
}

/*
Add invokes rtm.task.add to create a new task. This method uses Smart Add and can only add top-level tasks.
*/
func Add(apiToken string, name string, timelineID string) (task.Task, error) {
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
		return nil, err
	}

	return &response.Rsp.List.Taskseries[0], nil

}

type taskGetListResponse struct {
	Rsp struct {
		Stat  string
		Tasks struct {
			List []struct {
				ID         string
				Taskseries []taskResponse
			}
		}
	}
}

type GetListResponse struct {
	Rsp struct {
		Stat  string
		Tasks struct {
			List []struct {
				ID         string
				Taskseries []task.Task
			}
		}
	}
}

func GetList(apiToken string, filter string) (*GetListResponse, error) {
	args := map[string]string{
		"api_key":    api.APIKey,
		"auth_token": apiToken,
	}

	if filter != "" {
		args["filter"] = filter
	}

	var response taskGetListResponse
	unmarshal := func(b []byte) error {
		return json.Unmarshal(b, &response)
	}

	err := api.GetMethod("rtm.tasks.getList", args, unmarshal)
	if err != nil {
		return []task.Task{}, err
	}

	log.WithFields(log.Fields{
		"response": response,
	}).Debug("")

	log.WithFields(log.Fields{
		"records": len(response.Rsp.Tasks.List),
	}).Debug("Parsed responses.")

	var result GetListResponse
	result.Rsp.Stat = response.Rsp.Stat

	for idx, list := range response.Rsp.Tasks.List {
		result.Rsp.Tasks.List[idx].ID = response.Rsp.Tasks.List[idx].ID
		tasks := make([]task.Task, 0)

		for i := 0; i < len(list.Taskseries); i++ {
			tasks = append(tasks, list.Taskseries[i])
		}
		task.Sort(tasks)
		result.Rsp.Tasks.List[idx].Taskseries = tasks
	}

	return &result, nil

}

func Complete(apiToken, taskId, listID, timeline string)
