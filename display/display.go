package display

import (
	"fmt"

	"github.com/vikstrous/blinkythingy"
)

type Display interface {
	Flush([]blinkythingy.Color) error
}

type Factory interface {
	Create(string, blinkythingy.MapConfig) (Display, error)
}

type factory struct {
	creators map[string]Creator
}

type Creator func(blinkythingy.MapConfig) (Display, error)
type NamedCreator struct {
	Name    string
	Creator Creator
}

func NewFactory(creators []NamedCreator) Factory {
	f := factory{creators: map[string]Creator{}}
	for _, creator := range creators {
		f.register(creator.Name, creator.Creator)
	}
	return &f
}

func (f *factory) register(name string, creator Creator) {
	f.creators[name] = creator
}

func (f *factory) Create(name string, mapConfig blinkythingy.MapConfig) (Display, error) {
	creator, ok := f.creators[name]
	if !ok {
		return nil, fmt.Errorf("No display found with name: %s", name)
	}
	return creator(mapConfig)
}
