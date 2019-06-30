package task

import (
	"strings"

	"github.com/dishbreak/gomilk/api/tasks"
	"github.com/kyokomi/emoji"
)

// Tags are representation of tags for a task.
type Tags []string

// UnpackTags reads in the API response for tags and creates a display-friendly version.
func UnpackTags(input tasks.TaskResponse) Tags {
	tags := make([]string, len(input.Tags.Tag))

	for idx, val := range input.Tags.Tag {
		tags[idx] = emoji.Sprintf(":hash:%s", val)
	}

	return Tags(tags)
}

func (t Tags) String() string {
	return strings.Join(t, " ")
}
