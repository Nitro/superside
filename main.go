package main

import (
	"gopkg.in/alecthomas/kingpin.v1"
	"github.com/nitro/superside/tracker"
	"github.com/nitro/superside/persistence"
)

type CliOpts struct {
	ConfigFile *string
}

var state *tracker.Tracker

func parseCommandLine() *CliOpts {
	var opts CliOpts
	opts.ConfigFile = kingpin.Flag("config-file", "The config file to use").Short('f').Default("superside.toml").String()
	kingpin.Parse()
	return &opts
}

func main() {
	opts := parseCommandLine()
	config := parseConfig(*opts.ConfigFile)

	store := persistence.NewFileStore("data/")

	state = tracker.NewTracker(tracker.INITIAL_RING_SIZE, store)
	go state.ProcessUpdates()
	go state.ManagePersistence()

	serveHttp(config.Superside.BindIP, config.Superside.BindPort)
}
