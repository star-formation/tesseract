/*  Copyright 2019 The tesseract Authors

    This file is part of tesseract.

    tesseract is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    tesseract is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli"

	"github.com/star-formation/tesseract"
)

func init() {
	log.Root().SetHandler(log.MultiHandler(
		log.StreamHandler(os.Stderr, log.TerminalFormat(true)),
		log.LvlFilterHandler(
			log.LvlDebug,
			log.Must.FileHandler("starburst_errors.json", log.JSONFormat()))))
}

func main() {
	app := cli.NewApp()
	app.Name = "starburst"
	app.Version = "0.0.0"
	app.Usage = "Stretch to the Heavens"

	app.Flags = []cli.Flag{
		// TODO: override random beacon source and remove default value
		&cli.Uint64Flag{
			Name:  "testseed",
			Value: 0,
			Usage: "sets `UINT64` as deterministic seed",
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Info("==== starburst")
		// pass static seed to get deterministic, reproducible procgen
		tesseract.DevWorld(c.Uint64("testseed"))
		tesseract.StartWebSocket()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error("app.Run:", "err", err)
	}
}
