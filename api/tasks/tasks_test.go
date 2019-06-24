package tasks_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/dishbreak/gomilk/api/tasks"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalListResponse(t *testing.T) {
	filepath := path.Join("testdata", "list_response.json")

	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Error reading test data: %s\n", err)
	}

	var resp tasks.GetListResponse
	err = json.Unmarshal(buf, &resp)
	if err != nil {
		t.Fatalf("Error unmarshalling json: %s\n", err)
	}

	fmt.Printf("response: %s\n", resp)

	assert.Equal(t, 8, len(resp.Rsp.Tasks.List[0].Taskseries))

	names := make([]string, 8)

	for idx, taskSeries := range resp.Rsp.Tasks.List[0].Taskseries {
		names[idx] = taskSeries.RawName
	}

	expectedNames := []string{
		"Pack up donations",
		"Change sheets",
		"Look for software scale in LM boxes",
		"Activate Target redcard",
		"Search for gorillapod",
		"Install new insurance cards",
		"Sweep and mop floors",
		"Put donations out",
	}

	assert.Equal(t, expectedNames, names)
}
