package timelines

import (
	"encoding/json"

	"github.com/dishbreak/gomilk/api"
)

/*
TimelineCreateResponse contains the response to the call to rtm.timelines.create
*/
type TimelineCreateResponse struct {
	Rsp struct {
		Stat     string
		Timeline string
	}
}

/*
Create will invoke rtm.timelines.create. It will return the timeline ID as a string.

See details https://www.rememberthemilk.com/services/api/methods/rtm.timelines.create.rtm
*/
func Create(apiToken string) (*TimelineCreateResponse, error) {
	args := map[string]string{
		"auth_token": apiToken,
		"api_key":    api.APIKey,
	}

	var timelineResp TimelineCreateResponse
	unmarshall := func(b []byte) error {
		return json.Unmarshal(b, &timelineResp)
	}

	err := api.GetMethod("rtm.timelines.create", args, unmarshall)
	if err != nil {
		return nil, err
	}

	return &timelineResp, nil
}
