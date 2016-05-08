package github

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/colorutil"
	"github.com/vikstrous/blinkythingy/fetcher"
)

type GithubConfig struct {
	Project  string
	Username string
	Password string
	Query    string
}

func MapToGithubConfig(mapConfig blinkythingy.MapConfig) (GithubConfig, error) {
	config := GithubConfig{}
	marshalled, err := yaml.Marshal(mapConfig)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(marshalled, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

type githubFetcher struct {
	client   *github.Client
	query    string
	owner    string
	repo     string
	lock     sync.Mutex
	statuses []blinkythingy.Color
}

func New(mapConfig blinkythingy.MapConfig) (fetcher.Fetcher, error) {
	config, err := MapToGithubConfig(mapConfig)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(config.Project, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("github project name must contain a slash")
	}
	owner := parts[0]
	repo := parts[1]

	transport := github.BasicAuthTransport{
		Username: config.Username,
		Password: config.Password,
	}
	client := github.NewClient(transport.Client())

	return &githubFetcher{
		client: client,
		query:  fmt.Sprintf("is:pr repo:%s/%s %s is:open", owner, repo, config.Query),
		owner:  owner,
		repo:   repo,
	}, nil
}

func DebugRateLimit(res *github.Response) {
	remaining := res.Header.Get("X-RateLimit-Remaining")
	logrus.Debugf("Rate limit remaining: %s", remaining)
	reset := res.Header.Get("X-RateLimit-Reset")
	resetInt, err := strconv.ParseInt(reset, 10, 64)
	if err != nil {
		logrus.Debugf("Failed to parse X-RateLimit-Reset: %s", err)
		return
	}
	resetTime := time.Unix(resetInt, 0)
	logrus.Debugf("Rate limit reset: %s", resetTime)
}

func (f *githubFetcher) FetchStatuses() error {
	logrus.Debug("fetching from github")
	results, res, err := f.client.Search.Issues(f.query, nil)
	if err != nil {
		return err
	}
	DebugRateLimit(res)

	newStatuses := []string{}
	for _, issue := range results.Issues {
		id := issue.Number
		commits, res, err := f.client.PullRequests.ListCommits(f.owner, f.repo, *id, nil)
		if err != nil {
			return err
		}
		DebugRateLimit(res)
		sha := commits[len(commits)-1].SHA
		status, res, err := f.client.Repositories.GetCombinedStatus(f.owner, f.repo, *sha, nil)
		if err != nil {
			return err
		}
		DebugRateLimit(res)
		newStatuses = append(newStatuses, *status.State)
	}
	f.lock.Lock()
	defer f.lock.Unlock()
	f.statuses = GithubStatusesToColors(newStatuses)
	return nil
}

func (f *githubFetcher) ListStatuses() []blinkythingy.Color {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.statuses
}

func GithubStatusesToColors(statuses []string) []blinkythingy.Color {
	colors := []blinkythingy.Color{}
	for _, status := range statuses {
		colors = append(colors, GithubStatusToColor(status))
	}
	return colors
}

func GithubStatusToColor(status string) blinkythingy.Color {
	if status == "success" {
		return colorutil.MustParseColorString("green")
	} else if status == "pending" {
		color := colorutil.MustParseColorString("orange")
		color.BlinkPeriod = 10
		color.BlinkOn = 9
		return color
	} else if status == "failure" {
		return colorutil.MustParseColorString("red")
	} else {
		return colorutil.MustParseColorString("black")
	}
}
