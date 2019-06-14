package auth

import (
	"encoding/json"

	"github.com/dishbreak/gomilk/api"
)

// PermissionLevel describes the level you'd like to authenticate against.
type PermissionLevel int

const (
	// Read gives the ability to read task, contact, group and list details and contents.
	Read PermissionLevel = iota
	// Write gives the ability to add and modify task, contact, group and list details and contents (also allows you to read).
	Write
	// Delete gives the ability to delete tasks, contacts, groups and lists (also allows you to read and write).
	Delete
)

var permissionLevelStrings = [...]string{"read", "write", "delete"}

func (p PermissionLevel) String() string {
	return permissionLevelStrings[p]
}

type frob struct {
	Rsp struct {
		Stat string
		Frob string
	}
}

func GetFrob() (string, error) {
	args := map[string]string{
		"api_key": api.APIKey,
	}

	var frobResponse frob
	unmarshal := func(body []byte) error {
		return json.Unmarshal(body, &frobResponse)
	}

	err := api.GetMethod("rtm.auth.getFrob", args, unmarshal)
	if err != nil {
		return "", err
	}

	return frobResponse.Rsp.Frob, nil

}
