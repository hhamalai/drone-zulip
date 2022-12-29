package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
)

const TYPE_STREAM = "stream"
const TYPE_PRIVATE = "private"

const BuildStatusTemplate = `***{{.Emoji}} PIPELINE Report: {{.Status}} {{.Emoji}}***
***Pipeline:*** {{.Pipeline}}
***Repository / Branch / Commit:*** {{.Repository}} / {{.Branch}} / {{.Commit}}
***By:*** {{.By}}
***Link:*** {{.Link}}
***Message:*** {{.Message}}
`
const respFormat = `Webhook
  URL: %s
  RESPONSE STATUS: %s
  RESPONSE BODY: %s
`

type BuildStatus struct {
	Emoji      string
	Pipeline   string
	Status     string
	Commit     string
	Repository string
	Branch     string
	By         string
	Message    string
	Link       string
}

type (
	Repo struct {
		Owner string `json:"owner"`
		Name  string `json:"name"`
	}

	Build struct {
		Tag      string `json:"tag"`
		Event    string `json:"event"`
		Number   int    `json:"number"`
		Commit   string `json:"commit"`
		Ref      string `json:"ref"`
		Branch   string `json:"branch"`
		Author   string `json:"author"`
		Message  string `json:"message"`
		Status   string `json:"status"`
		Link     string `json:"link"`
		DeployTo string `json:"deployTo"`
		Started  int64  `json:"started"`
		Created  int64  `json:"created"`
	}

	Config struct {
		URL       string
		Type      string
		To        string
		Topic     string
		BotEmail  string
		BotApikey string
	}

	Stage struct {
		Type   string `json:"type"`
		Kind   string `json:"kind"`
		Name   string `json:"name"`
		Status string `json:"status"`
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
		Stage  Stage
	}
)

func (p *Plugin) Exec() error {
	uri, err := url.Parse(p.Config.URL)
	if err != nil {
		fmt.Printf("Error: Failed to parse the hook URL. %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Author", p.Build.Author)
	notification := BuildStatus{
		Pipeline:   p.Stage.Name,
		Commit:     p.Build.Commit,
		Repository: p.Repo.Name,
		Branch:     p.Build.Branch,
		By:         p.Build.Author,
		Message:    p.Build.Message,
		Link:       p.Build.Link,
	}

	if p.Config.Type != TYPE_STREAM && p.Config.Type != TYPE_PRIVATE {
		panic(errors.New(fmt.Sprintf("Type must be either %s or %s", TYPE_STREAM, TYPE_PRIVATE)))
	}

	if p.Stage.Status == "failure" || p.Build.Status == "failure" {
		notification.Status = "failure"
	} else if p.Stage.Status != "" {
		notification.Status = p.Stage.Status
	} else if p.Build.Status != "" {
		notification.Status = p.Build.Status
	} else {
		notification.Status = "unknown"
	}

	if notification.Status == "success" {
		notification.Emoji = ":check:"
	} else {
		notification.Emoji = ":fire:"
	}

	t, err := template.New("builds").Parse(BuildStatusTemplate)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, notification)
	if err != nil {
		panic(err)
	}

	values := url.Values{
		"type":    {p.Config.Type},
		"content": {buf.String()},
		"to":      {p.Config.To},
	}
	if p.Config.Type == TYPE_STREAM {
		if p.Config.Topic != "" {
			values["topic"] = []string{p.Config.Topic}
		} else {
			values["topic"] = []string{"no topic"}
		}
	}

	r := strings.NewReader(values.Encode())
	req, err := http.NewRequest(http.MethodPost, uri.String(), r)
	if err != nil {
		fmt.Printf("Error: Failed to create the HTTP request. %s\n", err)
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(p.Config.BotEmail, p.Config.BotApikey)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: Failed to execute the HTTP request. %s\n", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error: Failed to read the HTTP response body. %s\n", err)
		}
		fmt.Printf(respFormat, req.URL, resp.Status, string(body))
	}
	return err
}
