package api

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
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

	log.WithField("params", params).Debug("got arguments for request")

	query := u.Query()
	for key, value := range params {
		query.Set(key, value)
	}

	query.Set("api_sig", SignRequest(params))

	u.RawQuery = query.Encode()

	log.Debugf("creating URL: %s", u)

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

/*
TransactionRecord describes the transaction's identifiable information.

See https://www.rememberthemilk.com/services/api/timelines.rtm for details.
*/
type TransactionRecord struct {
	ID       string
	Undoable string
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

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		return err
	}

	log.Debugf("RTM API response: %s", string(prettyJSON.Bytes()))

	var errorMessage RTMAPIError
	if err := json.Unmarshal(body, &errorMessage); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to unmarshal api response.")
		return err
	}
	if errorMessage.Rsp.Stat == "fail" {
		return &errorMessage
	}

	if err := unmarshal(body); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to unmarshal api response.")
		return err
	}

	log.Debug("Completed API request.")

	return nil
}
