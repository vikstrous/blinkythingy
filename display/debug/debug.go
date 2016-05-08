package debug

import (
	"fmt"

	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/display"

	"github.com/aybabtme/rgbterm"
	"gopkg.in/yaml.v2"
)

type DebugConfig struct {
	Prefix string
}

func MapToDebugConfig(mapConfig blinkythingy.MapConfig) (DebugConfig, error) {
	config := DebugConfig{}
	marshalled, err := yaml.Marshal(mapConfig)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(marshalled, &config)
	if err != nil {
		return config, err
	}
	if config.Prefix == "" {
		config.Prefix = "Debug: "
	}
	return config, nil
}

type debugDisplay struct {
	prefix string
	tick   uint64
}

func New(mapConfig blinkythingy.MapConfig) (display.Display, error) {
	config, err := MapToDebugConfig(mapConfig)
	if err != nil {
		return nil, err
	}

	return &debugDisplay{
		prefix: config.Prefix,
	}, nil
}

func (d *debugDisplay) Flush(colors []blinkythingy.Color) error {
	fmt.Print(d.prefix)
	for _, color := range colors {
		if color.IsOn(d.tick) {
			fmt.Printf(rgbterm.FgString("X", color.R, color.G, color.B))
		} else {
			fmt.Printf(rgbterm.FgString("X", 0, 0, 0))
		}
	}
	fmt.Println("")
	d.tick++
	return nil
}
