package blinkytape

import (
	"io"
	"time"

	"github.com/tarm/serial"
)

type Color struct {
	R byte
	G byte
	B byte
}

func (c Color) WriteTo(writer io.Writer) error {
	if c.R == 255 {
		c.R = 254
	}
	if c.G == 255 {
		c.G = 254
	}
	if c.B == 255 {
		c.B = 254
	}
	_, err := writer.Write([]byte{c.R, c.G, c.B})
	return err
}

type BlinkyTape struct {
	serial   *serial.Port
	ledCount uint
}

func (b BlinkyTape) SendColors(colors []Color) error {
	colorCount := len(colors)
	for i, c := range colors {
		if i > int(b.ledCount) {
			break
		}
		err := c.WriteTo(b.serial)
		if err != nil {
			return err
		}
	}
	for i := 0; i < int(b.ledCount)-colorCount; i++ {
		err := Color{}.WriteTo(b.serial)
		if err != nil {
			return err
		}
	}
	// control character
	b.Send([]byte{255})
	return nil
}

func (b BlinkyTape) Send(data []byte) error {
	_, err := b.serial.Write(data)
	return err
}

func (b BlinkyTape) Close() error {
	return b.serial.Close()
}

func New(name string, ledCount uint) (*BlinkyTape, error) {
	config := serial.Config{
		Name:        name,
		Baud:        115200,
		ReadTimeout: time.Millisecond * 500,
	}

	port, err := serial.OpenPort(&config)
	if err != nil {
		return nil, err
	}

	blinky := BlinkyTape{
		serial:   port,
		ledCount: ledCount,
	}
	return &blinky, nil
}
