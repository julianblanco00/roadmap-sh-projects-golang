package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Event struct {
	Type  string `json:"type"`
	Actor struct {
		Id        int    `json:"id"`
		Login     string `json:"login"`
		AvatarUrl string `json:"avatar_url"`
	} `json:"actor"`
	Repo struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	CreatedAt string `json:"created_at"`
}

type EventData struct {
	Type      string
	Actor     string
	CreatedAt string
	Count     int
}

func fetchUserActivity() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\nEnter username: ")
		username, _ := reader.ReadString('\n')
		githubUrl := fmt.Sprintf("https://api.github.com/users/%s/events", strings.Trim(username, "\n"))

		resp, err := http.Get(githubUrl)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		if resp.StatusCode != 200 {
			fmt.Println("Error: ", resp.Status)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		data := []Event{}
		json.Unmarshal(body, &data)

		repoData := make(map[string]map[string]EventData)

		for _, event := range data {
			repoName := event.Repo.Name

			if _, ok := repoData[repoName]; !ok {
				repoData[repoName] = make(map[string]EventData)
			}

			currCount := repoData[repoName][event.Type].Count
			repoData[repoName][event.Type] = EventData{
				Type:      event.Type,
				Actor:     event.Actor.Login,
				CreatedAt: event.CreatedAt,
				Count:     currCount + 1,
			}
		}

		for repo, events := range repoData {
			for _, event := range events {
				fmt.Printf("%s %d times on %s\n", event.Type, event.Count, repo)
			}
		}
	}
}

func main() {
	fetchUserActivity()
}
