package main

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/newrelic/sidecar/catalog"
	"github.com/nitro/superside/circular"
	"github.com/nitro/superside/notification"
	"gopkg.in/alecthomas/kingpin.v1"
)

const (
	INITIAL_RING_SIZE   = 20
	CHANNEL_BUFFER_SIZE = 25
)

var (
	changes     *circular.Buffer
	changesChan chan catalog.StateChangedEvent
	listeners   []chan notification.Notification
	listenLock  sync.Mutex
)

type CliOpts struct {
	ConfigFile *string
}

func exitWithError(err error, message string) {
	if err != nil {
		log.Fatal("%s: %s", message, err.Error())
	}
}

func parseCommandLine() *CliOpts {
	var opts CliOpts
	opts.ConfigFile = kingpin.Flag("config-file", "The config file to use").Short('f').Default("superside.toml").String()
	kingpin.Parse()
	return &opts
}

// Subscribe a listener
func getListener() chan notification.Notification {
	listenChan := make(chan notification.Notification, 100)
	listenLock.Lock()
	listeners = append(listeners, listenChan)
	listenLock.Unlock()

	return listenChan
}

// Announce changes to all listeners
func tellListeners(evt *catalog.StateChangedEvent) {
	listenLock.Lock()
	defer listenLock.Unlock()

	// Try to tell the listener about the change but use a select
	// to protect us from any blocking readers.
	for _, listener := range listeners {
		select {
		case listener <- *notification.FromEvent(evt):
		default:
		}
	}
}

// Linearize the updates coming in from the async HTTP handler
func processUpdates() {
	for evt := range changesChan {
		changes.Insert(evt)
		tellListeners(&evt)
	}
}

func main() {
	opts := parseCommandLine()
	config := parseConfig(*opts.ConfigFile)

	changesChan = make(chan catalog.StateChangedEvent, CHANNEL_BUFFER_SIZE)
	changes = circular.NewBuffer(INITIAL_RING_SIZE)

	go processUpdates()

	serveHttp(config.Superside.BindIP, config.Superside.BindPort)
}
