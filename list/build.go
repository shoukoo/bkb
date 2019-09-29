package list

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/buildkite/go-buildkite/buildkite"
)

type Client struct {
	BKClient *buildkite.Client
	Builds   []Build
}

type Build struct {
	Pipeline     string
	Message      string
	Branch       string
	Status       string
	Commit       string
	Creator      string
	CreatorEmail string
	StartedAt    string
	ENV          string
	WebURL       string
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
	var builds []Build
	var limit = 1

	for {
		var p []buildkite.Build
		if r == nil {
			page = 1
		}

		listConfig := BuildListOptions(page)
		p, r, err = c.BKClient.Builds.ListByOrg(org, listConfig)
		builds = append(builds, reMapToBuild(p)...)

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

func reMapToBuild(b []buildkite.Build) []Build {
	var builds []Build
	for _, v := range b {
		build := Build{
			Message:      *v.Message,
			Branch:       *v.Branch,
			Pipeline:     *v.Pipeline.Slug,
			Status:       *v.State,
			Commit:       *v.Commit,
			Creator:      v.Creator.Name,
			CreatorEmail: v.Creator.Email,
			StartedAt:    v.StartedAt.Local().String(),
			ENV:          fmt.Sprintf("%v", v.Env),
			WebURL:       *v.WebURL,
		}

		builds = append(builds, build)

	}
	return builds
}

func (c *Client) Templates() *SelectTemplates {
	return &SelectTemplates{
		Active:   "> {{.Pipeline}} [ {{.Branch}} ]",
		Inactive: "  {{.Pipeline}} [ {{.Branch}} ]",
		Details: `
-------------- INFO --------------
Message: {{.Message}}
Branch: {{.Branch}}
Status: {{.Status}}
Commit: {{.Commit}}
Creator: {{.Creator}} ({{.CreatorEmail}})
Started: {{.StartedAt}}
ENV: {{.ENV}}
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
