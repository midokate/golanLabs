package main

import (
	"encoding/json"
	"io"
	"net/http"

	tea "charm.land/bubbletea/v2"
)

const (
	repoUrl = "https://api.github.com/orgs/charmbracelet/repos"
)

type repoGetSuccess []repo
type gotReposErrMsg error

type repo struct {
	Name string `json:"name"`
}

func getRepos() tea.Msg {
	req, err := http.NewRequest(http.MethodGet, repoUrl, nil)
	if err != nil {
		return nil
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "charmbracelet")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	defer func(Body io.ReadCloser) tea.Msg {
		err := Body.Close()
		if err != nil {
			return gotReposErrMsg(err)
		}
		return nil
	}(resp.Body)

	byteRepos, err := io.ReadAll(resp.Body)

	if err != nil {
		return gotReposErrMsg(err)
	}

	var repos []repo

	err = json.Unmarshal(byteRepos, &repos)

	if err != nil {
		return gotReposErrMsg(err)
	}

	return repoGetSuccess(repos)
}
