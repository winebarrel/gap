package main

import (
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/gap"
)

var (
	version string
)

func parseArgs() *gap.Options {
	var cli struct {
		gap.Options
		Version kong.VersionFlag
	}

	parser := kong.Must(&cli, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."
	_, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	return &cli.Options
}

func main() {
	options := parseArgs()
	server := gap.NewServer(options)
	err := server.Run()

	if err != nil {
		log.Fatal(err)
	}
}
