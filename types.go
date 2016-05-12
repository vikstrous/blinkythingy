package blinkythingy

import "fmt"

type Color struct {
	R uint8
	G uint8
	B uint8
	// number of ticks for length of period; 0==1, 1==2, etc.
	BlinkPeriod uint64
	// number of ticks to stay on during period; 0==1, 1==2, etc.
	BlinkOn uint64
}

func (c Color) String() string {
	return fmt.Sprintf("#%x", []byte{c.R, c.G, c.B})
}

func (c Color) IsOn(tick uint64) bool {
	blinkPeriod := c.BlinkPeriod + 1
	phase := tick % blinkPeriod
	if phase <= c.BlinkOn {
		return true
	}
	return false
}

type MapConfig map[string]interface{}

type Config struct {
	Debug      bool
	ReloadRate string
	BlinkRate  string
	Fetchers   []MapConfig
	Displays   []MapConfig
}

func (fc MapConfig) Type() string {
	typ, _ := fc["type"].(string)
	return typ
}

type HTTPClientConfig struct {
	InsecureTLS bool `yaml:"insecure-tls"`
	CA          string
}
