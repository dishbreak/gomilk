package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const (
	rootURL = "https://api.rememberthemilk.com/services/rest/"
)

/*
SignRequest uses a preshared secret to generate a signature that RTM will verify.

You should set the api_sig parameter with the return value.
*/
func SignRequest(params map[string]string) string {
	keys := make([]string, len(params))

	for key := range params {
		keys = append(keys, key)
	}
	pairs := make([]string, len(keys))
	sort.Strings(keys)

	for idx, key := range keys {
		pairs[idx] = key + params[key]
	}

	slug := []byte(SharedSecret + strings.Join(pairs, ""))

	return fmt.Sprintf("%x", md5.Sum(slug))
}

/*
FormURL creates a URL that you can use for making API requests or authenticating your users.
*/
func FormURL(baseURL string, method string, params map[string]string) *url.URL {
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("Failed to set root URL: %s\n", err)
	}

	params["method"] = method
	params["format"] = "json"

	query := u.Query()
	for key, value := range params {
		query.Set(key, value)
	}

	query.Set("api_sig", SignRequest(params))

	u.RawQuery = query.Encode()

	log.Printf("creating URL: %s", u)

	return u
}

var myClient http.Client

type RTMAPIError struct {
	Rsp struct {
		Stat string
		Err  struct {
			Code string
			Msg  string
		}
	}
}

func (e *RTMAPIError) Error() string {
	return fmt.Sprintf("RTM API encountered code %s: %s", e.Rsp.Err.Code, e.Rsp.Err.Msg)
}

/*
Get will issue a HTTP request using the GET verb.
*/
func GetMethod(method string, args map[string]string, unmarshal func([]byte) error) error {
	requestURL := FormURL(rootURL, method, args)

	resp, err := myClient.Get(requestURL.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(body))

	var errorMessage RTMAPIError
	if err := json.Unmarshal(body, &errorMessage); err != nil {
		return err
	}
	if errorMessage.Rsp.Stat == "fail" {
		return &errorMessage
	}

	if err := unmarshal(body); err != nil {
		return err
	}

	return nil
}
