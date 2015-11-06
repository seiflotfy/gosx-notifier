package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	osx "github.com/seiflotfy/github-notifier/osx"
)

var eventMapping = map[string]string{
	"WatchEvent":  "starred %s",
	"CreateEvent": "created %s",
	"PublicEvent": "made %s public",
	"IssuesEvent": "%s Issue %s/issues/%d",
	"ForkEvent":   "forked %s to %s",
	"PushEvent":   "pushed to %s at %s",
	"Default":     "%s %s",
}

type payload struct {
	Action  string                 `json:"action"`
	Issue   map[string]interface{} `json:"issue"`
	Forkee  map[string]interface{} `json:"forkee"`
	Payload map[string]interface{} `json:"payload"`
}

func formatIssuesMessage(val string, event github.Event) (string, string, error) {
	var p payload
	err := json.Unmarshal(*event.RawPayload, &p)
	if err != nil {
		fmt.Println(err)
	}
	var issueEvent github.IssueEvent
	err = json.Unmarshal(*event.RawPayload, &issueEvent)
	if err != nil {
		fmt.Println(err)
	}
	link := *issueEvent.Issue.HTMLURL
	message := fmt.Sprintf(val, p.Action, *event.Repo.Name, *issueEvent.Issue.Number)
	return message, link, nil
}

func formatForkEvent(val string, event github.Event) (string, string, error) {
	var p payload
	err := json.Unmarshal(*event.RawPayload, &p)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(event)
	link := p.Forkee["html_url"].(string)
	message := fmt.Sprintf(val, *event.Repo.Name, p.Forkee["full_name"].(string))
	return message, link, nil
}

func formatPushEvent(val string, event github.Event) (string, string, error) {
	var p payload
	err := json.Unmarshal(*event.RawPayload, &p)
	if err != nil {
		fmt.Println(err)
	}
	var pushEvent github.PushEvent
	err = json.Unmarshal(*event.RawPayload, &pushEvent)
	if err != nil {
		fmt.Println(err)
	}
	link := fmt.Sprintf("http://github.com/%s", *event.Repo.Name)
	message := fmt.Sprintf(val, *pushEvent.Ref, *event.Repo.Name)
	return message, link, nil
}

func formatDefaultEvent(val string, event github.Event) (string, string, error) {
	message := fmt.Sprintf(val, *event.Repo.Name)
	link := fmt.Sprintf("http://github.com/%s", *event.Repo.Name)
	return message, link, nil
}

func getEvents(username, password string) ([]github.Event, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/received_events?page=0", username)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Printf("Error : %s", err)
		return nil, err
	}

	req.SetBasicAuth(username, password)
	tc := &http.Client{}

	resp, err := tc.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
		return nil, err
	}

	// read body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error : %s", err)
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		fmt.Printf("Error : %s", err)
		return nil, err
	}

	var events []github.Event

	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Printf("Error : %s", err)
		return nil, err
	}
	return events, nil
}

func notify(event github.Event) {

	val, ok := eventMapping[*event.Type]
	if !ok {
		val = eventMapping["Default"]
	}

	var format func(string, github.Event) (string, string, error)

	switch {
	case *event.Type == "IssuesEvent":
		format = formatIssuesMessage
	case *event.Type == "ForkEvent":
		format = formatForkEvent
	case *event.Type == "PushEvent":
		format = formatPushEvent
	default:
		format = formatDefaultEvent
	}

	message, link, _ := format(val, event)
	note := osx.NewNotification(message)
	note.Link = link
	note.Title = *event.Actor.Login

	//Optionally, set a group which ensures only one notification is ever shown replacing previous notification of same group id.
	//note.Group = "Github-" + *event.ID
	//note.Sender = "com.geekyogre.github-notifier"
	//note.ContentImage = "GitHub-Mark-64px.png"
	//Optionally, set a sound from a predefined set.
	//note.Sound = gosxnotifier.Basso

	err := note.Push()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	username := ""
	password := ""
	ids := make(map[string]uint8)
	flag.StringVar(&username, "u", "", "Github username")
	flag.StringVar(&password, "p", "", "Github password")
	flag.Parse()

	for {
		events, err := getEvents(username, password)
		if err != nil {
			continue
		}

		nowM1 := time.Now().Add(-1 * time.Hour)

		for i, event := range events {
			if i == 5 || nowM1.Sub(event.CreatedAt.UTC()) > 0 {
				break
			}
			if _, ok := ids[*event.ID]; ok {
				continue
			}
			ids[*event.ID] = 0
			notify(event)
		}
		time.Sleep(30 * time.Second)
	}
}
