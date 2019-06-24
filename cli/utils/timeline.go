package utils

import (
	"github.com/dishbreak/gomilk/api/timelines"
	"github.com/dishbreak/gomilk/cli/login"
)

/*
Timeline will get the timeline that gomilk is using. If gomilk hasn't requested a timeline yet,
this call will request a new timeline from RTM.

For more info, see: https://www.rememberthemilk.com/services/api/timelines.rtm
*/
func Timeline() (string, error) {
	currentTimelineRecord, err := NewCurrentRecord("timeline")
	if err != nil {
		return "", err
	}

	// Attempt to use the timeline from the current record.
	// If it doesn't exist, request a new timeline from RTM
	timeline, err := currentTimelineRecord.Get()
	if err != nil {
		timelineRecord, err := timelines.Create(login.Token)
		timeline = timelineRecord.Rsp.Timeline
		if err != nil {
			return "", err
		}

		// persist the new timeline.
		currentTimelineRecord.Set(timeline)
	}

	return timeline, err

}
