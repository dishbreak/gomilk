package task

import (
	"encoding/json"
	"io/ioutil"

	"github.com/dishbreak/gomilk/api/timelines"

	"github.com/dishbreak/gomilk/cli/utils"

	"github.com/dishbreak/gomilk/api/tasks"
)

/*
Client provides a simplified interface that handles packing messages and obeying API semantics.
*/
type Client interface {
	Add(task string) (Task, error)
	List(filter string) ([]Task, error)
}

type client struct {
	token         string
	timelineID    string
	cacheFilePath string
}

/*
NewClient instantiates a new client object.
*/
func NewClient(token string) (Client, error) {
	cacheFilePath, err := utils.NewCache("tasks")
	if err != nil {
		return nil, err
	}

	timeline, err := timelines.Create(token)
	if err != nil {
		return nil, err
	}

	return &client{token: token, timelineID: timeline.Rsp.Timeline, cacheFilePath: cacheFilePath.Filename()}, nil
}

func (c *client) Add(task string) (Task, error) {
	resp, err := tasks.Add(c.token, task, c.timelineID)
	if err != nil {
		return nil, err
	}

	return unpackList(resp.Rsp.List)[0], nil
}

func (c *client) List(filter string) ([]Task, error) {
	resp, err := tasks.GetList(c.token, filter)
	if err != nil {
		return nil, err
	}

	result := make([]Task, 0)
	for _, tasklist := range resp.Rsp.Tasks.List {
		result = append(result, unpackList(tasklist)...)
	}

	Sort(result)

	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(c.cacheFilePath, data, 0600)
	if err != nil {
		return result, err
	}

	return result, nil

}

func unpackList(list tasks.TaskList) []Task {
	result := make([]Task, 0)
	for _, taskseries := range list.Taskseries {
		for _, taskInstance := range taskseries.Task {
			record := taskRecord{
				RawID:        taskInstance.RawID,
				RawName:      taskseries.RawName,
				TaskseriesID: taskseries.RawID,
				ListID:       list.ID,
				Due:          taskInstance.Due,
				HasDueTime:   taskInstance.HasDueTime,
			}

			result = append(result, record)
		}
	}

	return result
}
