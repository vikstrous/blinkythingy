package main

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/display"
	"github.com/vikstrous/blinkythingy/display/blinky"
	combineddisplay "github.com/vikstrous/blinkythingy/display/combined"
	"github.com/vikstrous/blinkythingy/display/debug"
	"github.com/vikstrous/blinkythingy/fetcher"
	"github.com/vikstrous/blinkythingy/fetcher/blank"
	combinedfetcher "github.com/vikstrous/blinkythingy/fetcher/combined"
	githubfetcher "github.com/vikstrous/blinkythingy/fetcher/github"
	"github.com/vikstrous/blinkythingy/fetcher/jenkins"
	"gopkg.in/yaml.v2"
)

type Flags struct {
	ConfigPath string
	Debug      bool
}

func RunFetcher(fetcher fetcher.Fetcher, display display.Display, reloadRate, blinkRate time.Duration) error {
	// we do one inline before we kick off the goroutine to make sure it works
	logrus.Debug("initial fetch")
	err := fetcher.FetchStatuses()
	if err != nil {
		return err
	}
	go func() {
		for {
			logrus.Debug("fetch")
			time.Sleep(reloadRate)
			err := fetcher.FetchStatuses()
			if err != nil {
				logrus.Warn(err)
			}
		}
	}()
	for {
		logrus.Debug("flush")
		colors := fetcher.ListStatuses()
		err := display.Flush(colors)
		if err != nil {
			logrus.Warn(err)
		}
		time.Sleep(blinkRate)
	}
}

func Run(config blinkythingy.Config) error {
	fetcherFactory := fetcher.NewFactory([]fetcher.NamedCreator{
		{"github", githubfetcher.New},
		{"blank", blank.New},
		{"jenkins", jenkins.New},
	})

	fetchers := []fetcher.Fetcher{}
	for _, f := range config.Fetchers {
		newFetcher, err := fetcherFactory.Create(f.Type(), f)
		if err != nil {
			return err
		}
		fetchers = append(fetchers, newFetcher)
	}

	displayFactory := display.NewFactory([]display.NamedCreator{
		{"blinky", blinky.New},
		{"debug", debug.New},
	})

	displays := []display.Display{}
	for _, d := range config.Displays {
		newDisplay, err := displayFactory.Create(d.Type(), d)
		if err != nil {
			return err
		}
		displays = append(displays, newDisplay)
	}

	if config.ReloadRate == "" {
		config.ReloadRate = "1m"
	}
	reloadRate, err := time.ParseDuration(config.ReloadRate)
	if err != nil {
		return err
	}
	if config.BlinkRate == "" {
		config.BlinkRate = "100ms"
	}
	blinkRate, err := time.ParseDuration(config.BlinkRate)
	if err != nil {
		return err
	}
	return RunFetcher(combinedfetcher.New(fetchers...), combineddisplay.New(displays...), reloadRate, blinkRate)
}

func main() {
	app := cli.NewApp()
	app.Name = "blinkytape"
	app.Usage = "change the color of your blinkytape"
	flags := Flags{}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			EnvVar:      "BLINKY_DEBUG",
			Destination: &flags.Debug,
		},
		cli.StringFlag{
			Name:        "config-path",
			EnvVar:      "BLINKY_CONFIG_PATH",
			Destination: &flags.ConfigPath,
			Value:       "config.yml",
		},
	}
	app.Action = func(c *cli.Context) {
		if flags.Debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		configBytes, err := ioutil.ReadFile(flags.ConfigPath)
		if err != nil {
			logrus.Fatal(err)
		}
		config := blinkythingy.Config{}
		err = yaml.Unmarshal(configBytes, &config)
		if err != nil {
			logrus.Fatal(err)
		}
		if config.Debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		err = Run(config)
		if err != nil {
			logrus.Fatal(err)
		}
	}
	//app.Commands = []cli.Command{}
	app.Run(os.Args)
}
