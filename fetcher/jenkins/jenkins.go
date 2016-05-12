package jenkins

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/colorutil"
	"github.com/vikstrous/blinkythingy/fetcher"
	"github.com/vikstrous/blinkythingy/jenkinsapi"
	"github.com/vikstrous/blinkythingy/util"
	"gopkg.in/yaml.v2"
)

type JenkinsConfig struct {
	blinkythingy.HTTPClientConfig
	Host     string
	Username string
	Password string
	Job      string
	Matrix   bool
}

func MapToJenkinsConfig(mapConfig blinkythingy.MapConfig) (JenkinsConfig, error) {
	config := JenkinsConfig{}
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

type jenkinsFetcher struct {
	client   jenkinsapi.Client
	job      string
	matrix   bool
	lock     sync.Mutex
	statuses []blinkythingy.Color
}

func New(mapConfig blinkythingy.MapConfig) (fetcher.Fetcher, error) {
	config, err := MapToJenkinsConfig(mapConfig)
	if err != nil {
		return nil, err
	}
	httpClient, err := util.HTTPClient(config.InsecureTLS, config.CA)
	if err != nil {
		return nil, err
	}
	jenkinsAPI := jenkinsapi.New(config.Host, config.Username, config.Password, httpClient)

	return &jenkinsFetcher{
		client: jenkinsAPI,
		job:    config.Job,
		matrix: config.Matrix,
	}, nil
}

func (jf *jenkinsFetcher) FetchStatuses() error {
	logrus.Debug("fetching from jenkins")
	statuses := []blinkythingy.Color{}
	// this is totally unsafe; TODO: catch panics here
	if jf.matrix {
		data, err := jf.client.JobDataWithFilter(jf.job, "activeConfigurations[color,name]")
		if err != nil {
			return err
		}
		dataMap := data.(map[string]interface{})
		configs := dataMap["activeConfigurations"]
		configsList := configs.([]interface{})
		for _, config := range configsList {
			configMap := config.(map[string]interface{})
			name := configMap["name"].(string)
			color := configMap["color"].(string)
			logrus.Debugf("Config: %s is %s", name, color)
			statuses = append(statuses, JenkinsColorToColor(color))
		}
	} else {
		data, err := jf.client.JobDataWithFilter(jf.job, "color")
		if err != nil {
			return err
		}
		dataMap := data.(map[string]interface{})
		color := dataMap["color"]
		statuses = append(statuses, JenkinsColorToColor(color.(string)))
	}
	jf.lock.Lock()
	defer jf.lock.Unlock()
	jf.statuses = statuses
	return nil
}

func (jf *jenkinsFetcher) ListStatuses() []blinkythingy.Color {
	jf.lock.Lock()
	defer jf.lock.Unlock()
	return jf.statuses
}

var green = colorutil.MustParseColorString("green")
var red = colorutil.MustParseColorString("red")
var orange = colorutil.MustParseColorString("orange")
var black = colorutil.MustParseColorString("black")

func JenkinsColorToColor(jc string) blinkythingy.Color {
	greenAnime := green
	greenAnime.BlinkPeriod = 10
	greenAnime.BlinkOn = 9
	redAnime := red
	redAnime.BlinkPeriod = 10
	redAnime.BlinkOn = 9
	orangeAnime := orange
	orangeAnime.BlinkPeriod = 10
	orangeAnime.BlinkOn = 9
	switch jc {
	case "blue":
		return green
	case "blue_anime":
		return greenAnime
	case "red":
		return red
	case "red_anime":
		return redAnime
	case "orange":
		return orange
	case "orange_anime":
		return orangeAnime
	case "notbuilt_anime":
		return orangeAnime
	}
	logrus.Debugf("Unknown jenkins status: %s", jc)
	return black
}
