package blank

import (
	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/fetcher"
	"gopkg.in/yaml.v2"
)

type BlankConfig struct {
	Number int
}

func MapToBlankConfig(mapConfig blinkythingy.MapConfig) (BlankConfig, error) {
	config := BlankConfig{}
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

type blankFetcher struct {
	blanks int
}

func New(mapConfig blinkythingy.MapConfig) (fetcher.Fetcher, error) {
	config, err := MapToBlankConfig(mapConfig)
	if err != nil {
		return nil, err
	}
	return &blankFetcher{config.Number}, nil
}

func (f *blankFetcher) FetchStatuses() error {
	return nil
}

func (f *blankFetcher) ListStatuses() []blinkythingy.Color {
	return Blank(f.blanks)
}

func Blank(num int) []blinkythingy.Color {
	colors := []blinkythingy.Color{}
	for i := 0; i < num; i++ {
		colors = append(colors, blinkythingy.Color{})
	}
	return colors
}
