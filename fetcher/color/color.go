package color

import (
	"fmt"

	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/colorutil"
	"github.com/vikstrous/blinkythingy/fetcher"
	"gopkg.in/yaml.v2"
)

type ColorConfig struct {
	Blink  bool
	Color  string
	Number int
}

func MapToColorConfig(mapConfig blinkythingy.MapConfig) (ColorConfig, error) {
	config := ColorConfig{}
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

type colorFetcher struct {
	color  blinkythingy.Color
	colors int
}

func New(mapConfig blinkythingy.MapConfig) (fetcher.Fetcher, error) {
	config, err := MapToColorConfig(mapConfig)
	if err != nil {
		return nil, err
	}
	color := config.Color
	if color == "" {
		color = "black"
	}
	parsedColor := colorutil.ParseColorString(color)
	if parsedColor == nil {
		return nil, fmt.Errorf("Failed to parse color: %s", color)
	}
	if config.Blink {
		parsedColor.BlinkPeriod = 9
		parsedColor.BlinkOn = 8
	}
	return &colorFetcher{*parsedColor, config.Number}, nil
}

func (f *colorFetcher) FetchStatuses() error {
	return nil
}

func (f *colorFetcher) ListStatuses() []blinkythingy.Color {
	colors := []blinkythingy.Color{}
	for i := 0; i < f.colors; i++ {
		colors = append(colors, f.color)
	}
	return colors
}
