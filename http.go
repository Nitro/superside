package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/Nitro/sidecar/catalog"
	"github.com/Nitro/superside/tracker"
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
	Message        string
	ClusterLatches *tracker.ClusterEventsLatch
}

// The health check endpoint.
func healthHandler(response http.ResponseWriter, req *http.Request, _ httprouter.Params, state *tracker.Tracker) {
	defer req.Body.Close()
	response.Header().Set("Content-Type", "application/json")

	//errors := make([]string, 0)

	message, _ := json.Marshal(ApiStatus{
		Message: "Healthy!",
		ClusterLatches: state.EventsLatch,
	})

	response.Write(message)
}

// Returns the currently stored state as a JSON blob
func servicesHandler(response http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	defer req.Body.Close()
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET")

	message, _ := json.Marshal(state.GetSvcEventsList())
	response.Write(message)
}

// Returns the currently stored state as a JSON blob
func deploymentsHandler(response http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	defer req.Body.Close()
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET")

	message, _ := json.Marshal(state.GetDeployments())
	response.Write(message)
}

// Receives POSTed state updates from Sidecar instances
func updateHandler(response http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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

	message, _ := json.Marshal(ApiMessage{"OK"})
	response.Write(message)
}

// Handle the listening endpoint websocket
func listenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func uiRedirectHandler(response http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	http.Redirect(response, req, "/ui/", 301)
}

func makeTrackerHandler(fn func(http.ResponseWriter, *http.Request,
	httprouter.Params, *tracker.Tracker)) httprouter.Handle {

	return func(response http.ResponseWriter, req *http.Request, params httprouter.Params) {
		fn(response, req, params, state)
	}
}

// Start the HTTP server and begin handling requests. This is a
// blocking call.
func serveHttp(listenIp string, listenPort int, state *tracker.Tracker) {
	listenStr := fmt.Sprintf("%s:%d", listenIp, listenPort)

	log.Infof("Starting up on %s", listenStr)

	router := httprouter.New()
	router.GET("/", uiRedirectHandler)
	router.POST("/api/update", updateHandler)
	router.GET("/api/state/services", servicesHandler)
	router.GET("/api/state/deployments", deploymentsHandler)
	router.GET("/health", makeTrackerHandler(healthHandler))
	router.GET("/listen", listenHandler)
	router.ServeFiles("/ui/*filepath", http.Dir("public/app"))

	http.Handle("/", handlers.LoggingHandler(os.Stdout, router))
	err := http.ListenAndServe(listenStr, nil)
	if err != nil {
		log.Fatalf("Can't start http server: %s", err.Error())
	}
}
