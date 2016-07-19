package main

import (
	"gopkg.in/alecthomas/kingpin.v1"
	"github.com/nitro/superside/tracker"
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

	state = tracker.NewTracker(tracker.INITIAL_RING_SIZE)
	go state.ProcessUpdates()

	serveHttp(config.Superside.BindIP, config.Superside.BindPort)
}
