package main

import (
	"encoding/hex"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/vikstrous/go-blinkytape"
)

func main() {
	app := cli.NewApp()
	app.Name = "blinkytape"
	app.Usage = "change the color of your blinkytape"
	app.Commands = []cli.Command{
		{
			Name:  "set",
			Usage: "set the color of the whole strip",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "color",
					EnvVar: "COLOR",
				},
			},
			Action: func(c *cli.Context) {
				blinky, err := blinkytape.New("/dev/ttyACM0", 60)
				if err != nil {
					logrus.Errorf("error opening port: %s", err)
					return
				}
				parsed := make([]byte, 3)
				_, err = hex.Decode(parsed, []byte(c.String("color")))
				pattern := []blinkytape.Color{}
				if err != nil {
					logrus.Errorf("failed to parse color: %s", err)
					return
				}
				for i := 0; i < 60; i++ {
					pattern = append(pattern, blinkytape.Color{parsed[0], parsed[1], parsed[2]})
				}
				err = blinky.SendColors(pattern)
				if err != nil {
					logrus.Errorf("error sending: %s", err)
					return
				}
			},
		},
	}
	app.Run(os.Args)
}
