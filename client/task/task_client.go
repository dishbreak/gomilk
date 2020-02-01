package task

import (
	"encoding/json"
	"io/ioutil"

	"github.com/dishbreak/gomilk/cli/utils"

	"github.com/dishbreak/gomilk/api/tasks"
)

/*
Client provides a simplified interface that handles packing messages and obeying API semantics.
*/
type Client interface {
	Add(task string) (Task, error)
	List(filter string) ([]Task, error)
	Complete(task Task) (Task, error)
	GetCachedTasks() ([]Task, error)
	SetDueDate(task Task, due string) (Task, error)
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

	timeline, err := utils.Timeline(token)
	if err != nil {
		return nil, err
	}

	return &client{token: token, timelineID: timeline, cacheFilePath: cacheFilePath.Filename()}, nil
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

func (c *client) GetCachedTasks() ([]Task, error) {
	cachedTasks := make([]taskRecord, 0)
	result := make([]Task, 0)

	buf, err := ioutil.ReadFile(c.cacheFilePath)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(buf, &cachedTasks)
	if err != nil {
		return result, err
	}

	for _, val := range cachedTasks {
		result = append(result, val)
	}

	return result, nil

}

func (c *client) Complete(task Task) (Task, error) {

	taskID, taskseriesID, listID := task.ID()

	resp, err := tasks.Complete(c.token, c.timelineID, listID, taskseriesID, taskID)
	if err != nil {
		return nil, err
	}

	return unpackList(resp.Rsp.List)[0], nil
}

func (c *client) SetDueDate(task Task, due string) (Task, error) {
	taskID, taskseriesID, listID := task.ID()

	resp, err := tasks.SetDueDate(c.token, c.timelineID, listID, taskseriesID, taskID, due)
	if err != nil {
		return nil, err
	}

	return unpackList(resp.Rsp.List)[0], nil
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
				RawPriority:  taskInstance.Priority,
				RawTags:      UnpackTags(taskseries),
				Completed:    taskInstance.Completed,
			}

			result = append(result, record)
		}
	}

	return result
}
