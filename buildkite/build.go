package buildkite

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/buildkite/go-buildkite/buildkite"
	"github.com/shoukoo/am2/list"
)

type Client struct {
	BKClient *buildkite.Client
	Builds   []buildkite.Build
}

func BuildkiteClient() (*Client, error) {
	token := os.Getenv("TOKEN")

	if token == "" {
		return nil, fmt.Errorf("TOKEN env variable is missing, please go https://buildkite.com/user/api-access-tokens to get a new token")
	}

	config, err := buildkite.NewTokenConfig(token, false)

	if err != nil {
		return nil, err
	}

	client := &Client{
		BKClient: buildkite.NewClient(config.Client()),
	}

	return client, nil

}

func (c *Client) GetRecentBuilds(org string) error {
	var r *buildkite.Response
	var page int
	var err error
	var builds []buildkite.Build
	var limit = 1

	for {
		var p []buildkite.Build
		if r == nil {
			page = 1
		}

		listConfig := BuildListOptions(page)
		p, r, err = c.BKClient.Builds.ListByOrg(org, listConfig)
		builds = append(builds, p...)

		if err != nil {
			return err
		}

		if r.LastPage == 0 {
			break
		}

		if page == limit {
			break
		}

		page = r.NextPage
	}

	c.Builds = builds

	return nil
}

func (c *Client) Templates() *list.SelectTemplates {
	return &list.SelectTemplates{
		Active:   "> {{.Pipeline.Slug}} [ {{.Branch}} ]",
		Inactive: "  {{.Pipeline.Slug}} [ {{.Branch}} ]",
		Details: `
-------------- Info --------------
Message: {{.Message}}
Branch: {{.Branch}}
Status: {{.State}}
Commit: {{.Commit}}
Creator: {{.Creator.Name}} ({{.Creator.Email}})
Started: {{.StartedAt}}
ENV: {{.Env}}
`,
	}
}

// BuildListOptions return list pagination config
func BuildListOptions(page int) *buildkite.BuildsListOptions {
	return &buildkite.BuildsListOptions{
		ListOptions: buildkite.ListOptions{
			Page:    page,
			PerPage: 100,
		},
	}
}

func Open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
