package blinky

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/Sirupsen/logrus"
	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/display"
	"github.com/vikstrous/go-blinkytape"
)

type BlinkyConfig struct {
	Path string
}

func MapToBlinkyConfig(mapConfig blinkythingy.MapConfig) (BlinkyConfig, error) {
	config := BlinkyConfig{}
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

type blinkyDisplay struct {
	blinky *blinkytape.BlinkyTape
	path   string
	tick   uint64
}

func New(mapConfig blinkythingy.MapConfig) (display.Display, error) {
	config, err := MapToBlinkyConfig(mapConfig)
	if err != nil {
		return nil, err
	}

	d := &blinkyDisplay{
		path: config.Path,
	}

	err = d.openConn()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *blinkyDisplay) openConn() error {
	path := d.path
	if path == "" {
		files, _ := ioutil.ReadDir("/dev")
		for _, f := range files {
			if strings.HasPrefix(f.Name(), "ttyACM") {
				path = filepath.Join("/dev", f.Name())
			}
		}
	}
	if path == "" {
		return fmt.Errorf("Failed to find blinky device")
	}
	// try to recover
	blinky, err := blinkytape.New(path, 60)
	if err != nil {
		return fmt.Errorf("error opening port at %s: %s", path, err)
	}
	if d.blinky != nil {
		err := d.blinky.Close()
		if err != nil {
			logrus.Warnf("error closing blinky device: %s", err)
		}
	}
	d.blinky = blinky
	return nil
}

func (d *blinkyDisplay) Flush(colors []blinkythingy.Color) error {
	defer func() { d.tick++ }()
	blinkyColors := []blinkytape.Color{}
	for _, color := range colors {
		if color.IsOn(d.tick) {
			blinkyColors = append(blinkyColors, blinkytape.Color{
				R: color.R,
				G: color.G,
				B: color.B,
			})
		} else {
			blinkyColors = append(blinkyColors, blinkytape.Color{})
		}
	}
	err := d.blinky.SendColors(blinkyColors)
	if err != nil {
		return d.openConn()
	}
	return nil
}
