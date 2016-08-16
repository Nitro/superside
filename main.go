package main

import (
	"gopkg.in/alecthomas/kingpin.v1"
	"github.com/nitro/superside/tracker"
	"github.com/nitro/superside/persistence"
)

type CliOpts struct {
	ConfigFile *string
	Persist    *bool
}

var state *tracker.Tracker

func parseCommandLine() *CliOpts {
	var opts CliOpts
	opts.ConfigFile = kingpin.Flag("config-file", "The config file to use").Short('f').Default("superside.toml").String()
	opts.Persist = kingpin.Flag("persist", "Do we persist and load data from the store?").Short('p').Default("true").Bool()
	kingpin.Parse()
	return &opts
}

func main() {
	opts := parseCommandLine()
	config := parseConfig(*opts.ConfigFile)

	var store persistence.Store
	if *opts.Persist {
		store = persistence.NewFileStore("data/")
	} else {
		store = &persistence.NoopStore{}
	}

	state = tracker.NewTracker(tracker.INITIAL_RING_SIZE, store)
	go state.ProcessUpdates()
	go state.ManagePersistence()

	serveHttp(config.Superside.BindIP, config.Superside.BindPort, state)
}
