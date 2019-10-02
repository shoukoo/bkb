package list

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/99designs/keyring"
	"github.com/buildkite/go-buildkite/buildkite"
)

type Client struct {
	BKClient *buildkite.Client
	Builds   []Build
	Org      string
}

type Build struct {
	Pipeline     string
	Message      string
	Branch       string
	Status       string
	Commit       string
	Creator      string
	CreatorEmail string
	CreatedAt    string
	ENV          string
	WebURL       string
	Elapsed      string
}

func BuildkiteClient() (*Client, error) {

	ring, err := keyring.Open(keyring.Config{
		ServiceName: "buildkite-beaver",
	})
	if err != nil {
		return nil, err
	}

	token, err := ring.Get("token")
	if err != nil {
		return nil, err
	}

	org, err := ring.Get("org")
	if err != nil {
		return nil, err
	}

	config, err := buildkite.NewTokenConfig(string(token.Data), false)
	if err != nil {
		return nil, err
	}

	client := &Client{
		BKClient: buildkite.NewClient(config.Client()),
		Org:      string(org.Data),
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
		elapsed := time.Since(v.CreatedAt.Time)
		created := v.CreatedAt.Local().Format("Monday, 2-January-2006")
		build := Build{
			Message:      *v.Message,
			Branch:       *v.Branch,
			Pipeline:     *v.Pipeline.Slug,
			Status:       *v.State,
			Commit:       *v.Commit,
			Creator:      v.Creator.Name,
			CreatorEmail: v.Creator.Email,
			CreatedAt:    created,
			Elapsed:      elapsed.Truncate(time.Second).String(),
			ENV:          fmt.Sprintf("%v", v.Env),
			WebURL:       *v.WebURL,
		}

		builds = append(builds, build)

	}
	return builds
}

func (c *Client) Templates() *SelectTemplates {
	return &SelectTemplates{
		Active:   `{{"▶" | cyan}} [ {{.Branch | cyan}} | {{.Status | cyan}} ] {{.Pipeline | cyan}} `,
		Inactive: `  [ {{.Branch}} | {{.Status}} ] {{.Pipeline}}`,
		Details: `
{{"---------------------------------------------" | blue}}
Message: {{.Message}}
Branch:  {{.Branch}}
Status:  {{.Status}}
Age:     {{.Elapsed}}
Commit:  {{.Commit}}
Creator: {{.Creator}} ({{.CreatorEmail}})
Started: {{.CreatedAt}}
ENV:     {{.ENV}}
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
