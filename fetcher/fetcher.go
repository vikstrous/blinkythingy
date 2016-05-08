package fetcher

import (
	"fmt"

	"github.com/vikstrous/blinkythingy"
)

type Fetcher interface {
	FetchStatuses() error
	ListStatuses() []blinkythingy.Color
}

type Factory interface {
	Create(string, blinkythingy.MapConfig) (Fetcher, error)
}

type factory struct {
	creators map[string]Creator
}

type Creator func(blinkythingy.MapConfig) (Fetcher, error)
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

func (f *factory) Create(name string, mapConfig blinkythingy.MapConfig) (Fetcher, error) {
	creator, ok := f.creators[name]
	if !ok {
		return nil, fmt.Errorf("No fetcher found with name: %s", name)
	}
	return creator(mapConfig)
}
