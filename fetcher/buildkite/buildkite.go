package buildkite

import (
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/colorutil"
	"github.com/vikstrous/blinkythingy/fetcher"
	"gopkg.in/buildkite/go-buildkite.v2/buildkite"
	"gopkg.in/yaml.v2"
)

type BuildkiteConfig struct {
	Token    string
	Org      string
	Pipeline string
	Branch   string
	Limit    int
	Debug    bool
}

func MapToBuildkiteConfig(mapConfig blinkythingy.MapConfig) (BuildkiteConfig, error) {
	config := BuildkiteConfig{}
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

type buildkiteFetcher struct {
	client   *buildkite.Client
	org      string
	pipeline string
	branch   string
	limit    int
	lock     sync.Mutex
	statuses []blinkythingy.Color
}

func New(mapConfig blinkythingy.MapConfig) (fetcher.Fetcher, error) {
	config, err := MapToBuildkiteConfig(mapConfig)
	if err != nil {
		return nil, err
	}

	// We don't support custom CAs and insecure TLS for now.
	// See: https://github.com/buildkite/go-buildkite/issues/14
	// We could work around it if necessary, but it's low priority.
	clientConfig, err := buildkite.NewTokenConfig(config.Token, config.Debug)
	if err != nil {
		return nil, fmt.Errorf("NewTokenConfig failed: %s", err)
	}

	client := buildkite.NewClient(clientConfig.Client())

	if config.Limit == 0 {
		config.Limit = 1
	}

	return &buildkiteFetcher{
		client:   client,
		org:      config.Org,
		pipeline: config.Pipeline,
		branch:   config.Branch,
		limit:    config.Limit,
	}, nil
}

func (b *buildkiteFetcher) FetchStatuses() error {
	logrus.Debug("fetching from buildkite")
	statuses := []blinkythingy.Color{}

	builds, resp, err := b.client.Builds.ListByPipeline(b.org, b.pipeline, &buildkite.BuildsListOptions{Branch: b.branch})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	count := 0
	for _, build := range builds {
		if build.State != nil {
			statuses = append(statuses, BuildkiteStateToColor(*build.State))
			count++
			if count >= b.limit {
				break
			}
		}
	}

	b.lock.Lock()
	defer b.lock.Unlock()
	b.statuses = statuses
	return nil
}

func (b *buildkiteFetcher) ListStatuses() []blinkythingy.Color {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.statuses
}

var green = colorutil.MustParseColorString("green")
var red = colorutil.MustParseColorString("red")
var orange = colorutil.MustParseColorString("orange")
var yellow = colorutil.MustParseColorString("yellow")
var blue = colorutil.MustParseColorString("blue")
var black = colorutil.MustParseColorString("black")

// possible values: running, scheduled, passed, failed, canceled, skipped and not_run
func BuildkiteStateToColor(status string) blinkythingy.Color {
	yellowAnime := yellow
	yellowAnime.BlinkPeriod = 10
	yellowAnime.BlinkOn = 9
	switch status {
	case "passed":
		return green
	case "failed":
		return red
	case "scheduled":
		return orange
	case "running":
		return yellowAnime
	case "canceled":
		return blue
	case "skipped":
		return blue
	case "not_run":
		return blue
	}
	logrus.Debugf("Unknown buildkite status: %s", status)
	return black
}
