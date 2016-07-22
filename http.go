package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/newrelic/sidecar/catalog"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 4096,
}

type ApiErrors struct {
	Errors []string
}

type ApiMessage struct {
	Message string
}

type ApiStatus struct {
	Message     string
	LastChanged time.Time
}

// The health check endpoint. Tells us if HAproxy is running and has
// been properly configured. Since this is critical infrastructure this
// helps make sure a host is not "down" by havign the proxy down.
func healthHandler(response http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	response.Header().Set("Content-Type", "application/json")

	//errors := make([]string, 0)

	message, _ := json.Marshal(ApiStatus{
		Message: "Healthy!",
	})

	response.Write(message)
}

// Returns the currently stored state as a JSON blob
func servicesHandler(response http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	response.Header().Set("Content-Type", "application/json")

	message, _ := json.Marshal(state.GetSvcEventsList())
	response.Write(message)
}

// Returns the currently stored state as a JSON blob
func deploymentsHandler(response http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	response.Header().Set("Content-Type", "application/json")

	message, _ := json.Marshal(state.GetDeployments())
	response.Write(message)
}

// Receives POSTed state updates from Sidecar instances
func updateHandler(response http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	response.Header().Set("Content-Type", "application/json")

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		message, _ := json.Marshal(ApiErrors{[]string{err.Error()}})
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(message)
		return
	}

	var evt catalog.StateChangedEvent
	err = json.Unmarshal(data, &evt)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	state.EnqueueUpdate(evt) // Potentially blocking

	message, _ := json.Marshal(ApiMessage{Message: "OK"})
	response.Write(message)
}

// Handle the listening endpoint websocket
func listenHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}

	svcEventsChan := state.GetSvcEventsListener()
	defer state.RemoveSvcEventsListener(svcEventsChan)

	deployChan := state.GetDeploymentListener()
	defer state.RemoveDeploymentListener(deployChan)

	// Loop, multiplexing the two channels and constructing events
	// from each.
	for {
		var message []byte

		select {
		case evt := <-svcEventsChan:
			output := struct {
				Type string
				Data interface{}
			}{"ServiceEvent", evt}
			message, err = json.Marshal(output)

		case deploy := <-deployChan:
			output := struct {
				Type string
				Data interface{}
			}{"Deployment", deploy}
			message, err = json.Marshal(output)
		}

		if err != nil {
			log.Error("Error marshaling JSON event " + err.Error())
			continue
		}

		if err = conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Warn(err.Error())
			return
		}
	}
}

// Start the HTTP server and begin handling requests. This is a
// blocking call.
func serveHttp(listenIp string, listenPort int) {
	listenStr := fmt.Sprintf("%s:%d", listenIp, listenPort)

	log.Infof("Starting up on %s", listenStr)
	fs := http.FileServer(http.Dir("public/app/"))
	router := mux.NewRouter()

	router.HandleFunc("/api/update", updateHandler).Methods("POST")
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/state/services", servicesHandler).Methods("GET")
	router.HandleFunc("/api/state/deployments", deploymentsHandler).Methods("GET")
	router.HandleFunc("/listen", listenHandler).Methods("GET")
	router.PathPrefix("/").Handler(fs)
	http.Handle("/", handlers.LoggingHandler(os.Stdout, router))

	err := http.ListenAndServe(listenStr, nil)
	if err != nil {
		log.Fatalf("Can't start http server: %s", err.Error())
	}
}
