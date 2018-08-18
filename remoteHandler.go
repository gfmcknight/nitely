package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

// Clones down a repository based on its information
func cloneOrFetch(build buildInfo) {
	os.Chdir(getStorageBase())
	if _, err := os.Stat(build.AbsolutePath); os.IsNotExist(err) {
		err := os.MkdirAll(build.AbsolutePath, os.ModeDir)
		if err != nil {
			fmt.Printf("Error creating the path %s", build.AbsolutePath)
		}

		// TODO: Fix (potential security risk)
		url := fmt.Sprintf("https://%s@github.com/%s.git", getToken(), build.Remote)

		cmd := exec.Command("git", "clone", url, build.AbsolutePath)
		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
			fmt.Printf("... When trying to clone %s\n", url)
		}
	}

	// Fetch updates from the remote if necessary
	os.Chdir(build.AbsolutePath)
	cmd := exec.Command("git", "fetch")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("... When trying to fetch for %s\n", build.Name)
	}

}
