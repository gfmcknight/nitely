package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path"
	"time"
)

// Access token for github can be stored in a file in plaintext
// or as a property
func getToken() string {

	tokenPtr := getProperty(nil, "github.token")

	if tokenPtr != nil {
		return *tokenPtr
	}

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

// Steals a final design philosophy statement from Github
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

// Makes a github api request to set the build status
// for a given commit
func setStatus(repoName, status, commitID string) {
	queries := make(map[string]interface{})
	body := make(map[string]interface{})
	body["state"] = status
	body["taget_url"] = nil
	body["description"] = "The build ended in " + status
	body["context"] = "continuous-integration/nitely"

	url := fmt.Sprintf("/repos/%s/statuses/%s", repoName, commitID)
	resp, err := requestJSON("POST", url, getToken(), queries, body)
	if resp.StatusCode != 201 {
		rb, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(rb))
	}

	if err != nil {
		fmt.Println(err)
	}
}

func getCommitID() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	return string(out)[0:40]
}

// Creates a remote from a specified repository. The remote
// comes in the form user/repository.
func addFromRemote(name, remote, branch string) {
	if name == "" || branch == "" {
		panic("Remotes specify name, and branch!")
	}

	info := buildInfo{}
	info.Name = name
	info.Branch = branch
	info.Remote = remote
	info.AbsolutePath = path.Join(getStorageBase(), "remotes", remote)

	fmt.Printf("Added build %s on branch %s from repo %s\n", info.Name, info.Branch, info.Remote)

	insertBuildInfo(nil, info)
}
