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

/*
FrobRecord contains a Golang representation of the response to the call rtm.auth.getFrob
*/
type FrobRecord struct {
	Rsp struct {
		Stat string
		Frob string
	}
}

/*
GetFrob calls rtm.auth.getFrob using the preconfigured API key.
https://www.rememberthemilk.com/services/api/methods/rtm.auth.getFrob.rtm
*/
func GetFrob() (*FrobRecord, error) {
	args := map[string]string{
		"api_key": api.APIKey,
	}

	var frobResponse FrobRecord
	unmarshal := func(body []byte) error {
		return json.Unmarshal(body, &frobResponse)
	}

	rawResponse, err := api.GetMethod("rtm.auth.getFrob", args, unmarshal)
	if err != nil {
		return nil, err
	}
	frobResponse, _ = rawResponse.(FrobRecord)

	return &frobResponse, nil

}

/*
TokenRecord contains a Golang representation of the response to the call rtm.auth.getToken
*/
type TokenRecord struct {
	Rsp struct {
		Stat string
		Auth *TokenCheckRecord
	}
}

/*
TokenCheckRecord contains a Golang representation of the response to the call rtm.auth.checkToken
*/
type TokenCheckRecord struct {
	Token string
	Perms string
	User  struct {
		ID       string
		Username string
		Fullname string
	}
}

/*
GetToken calls rtm.auth.getToken using the preconfigured API key.
https://www.rememberthemilk.com/services/api/methods/rtm.auth.getToken.rtm

GetToken requires a frob from the RTM API, use GetFrob() to get one.
*/
func GetToken(frob string) (*TokenRecord, error) {
	args := map[string]string{
		"frob":    frob,
		"api_key": api.APIKey,
	}

	var tokenRecord TokenRecord

	rawResponse, err := api.GetMethod("rtm.auth.getToken", args, tokenRecord)
	if err != nil {
		return nil, err
	}
	tokenRecord, _ = rawResponse.(TokenRecord)
	return &tokenRecord, nil
}

/*
CheckToken calls rtm.auth.checkToken using the preconfigured API key.
https://www.rememberthemilk.com/services/api/methods/rtm.auth.checkToken.rtm

You'll need to get token using GetToken() first.
*/
func CheckToken(authToken string) (*TokenCheckRecord, error) {
	args := map[string]string{
		"auth_token": authToken,
		"api_key":    api.APIKey,
	}

	var tokenCheckRecord TokenCheckRecord

	rawResponse, err := api.GetMethod("rtm.auth.checkToken", args, tokenCheckRecord)
	if err != nil {
		return nil, err
	}

	tokenCheckRecord, _ = rawResponse.(TokenCheckRecord)

	return &tokenCheckRecord, nil
}

const (
	//ErrorInvalidAuthToken represents the situatuon where the token is invalidated (expired or access revoked)
	ErrorInvalidAuthToken = "98"
)
