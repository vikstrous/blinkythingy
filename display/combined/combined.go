package combined

import (
	"fmt"

	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/display"
)

type combinedDisplay struct {
	displays []display.Display
}

func New(displays ...display.Display) display.Display {
	return &combinedDisplay{
		displays: displays,
	}
}

func (d *combinedDisplay) Flush(colors []blinkythingy.Color) error {
	errs := []error{}
	for _, display := range d.displays {
		err := display.Flush(colors)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", errs)
	}
	return nil
}
