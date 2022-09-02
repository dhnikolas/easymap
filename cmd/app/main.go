package main

import (
	"log"
	"os"

	"github.com/dhnikolas/easymap/internal/command"
	cli "github.com/urfave/cli/v2"
)

func main() {

	easyMapApp := &command.Application{}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "gen",
				Aliases: []string{"g"},
				Usage:   "generate mapping for struct A to struct B (example: easymap gen /path/to/file/main.go:SomeStruct /path/to/file/item.go:SimpleStruct)",
				Action:  easyMapApp.Gen,
			},
			{
				Name:    "copy",
				Aliases: []string{"c"},
				Usage:   "copy struct (example: easymap copy /path/to/file/main.go:SomeStruct NewStructName)",
				Action:  easyMapApp.Copy,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
