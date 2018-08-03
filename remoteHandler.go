package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

func getToken() string {
	tokenFile := path.Join(getStorageBase(), ".token")
	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(token)
}

func request(method, url, token, queryString string,
	body []byte) (*http.Response, error) {

	url = "https://api.github.com" + url + queryString

	client := http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}

	return client.Do(req)
}

func requestString(method, url, token, queryString, body string) (*http.Response, error) {
	return request(method, url, token, queryString, []byte(body))
}

func requestJSON(method, url, token string,
	queries, body map[string]interface{}) (*http.Response, error) {

	queryString := ""
	if len(queries) > 0 {
		queryString = "?"
	}
	for k, v := range queries {
		queryString = fmt.Sprintf("%s%s=%v&", queryString, k, v)
	}

	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return request(method, url, token, queryString, json)
}

func getZen() string {
	resp, err := requestString("GET", "/zen", getToken(), "", "")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(body)
}

func setStatus(owner, repoName, status, commitID string) {
	queries := make(map[string]interface{})
	body := make(map[string]interface{})
	body["state"] = status
	body["taget_url"] = nil
	body["description"] = "The build ended in " + status
	body["context"] = "continuous-integration/nitely"

	url := fmt.Sprintf("/repos/%s/%s/statuses/%s", owner, repoName, commitID)
	requestJSON("POST", url, getToken(), queries, body)
}
