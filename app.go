package gapp

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Hearst-DD/gappconfig"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/stretchr/graceful"
)

// HandlerMapping is used to allow an app to attach its handlers to the Gorilla mux.
type HandlerMapping struct {
	Route   string
	Handler func(rw http.ResponseWriter, r *http.Request)
}

// Gapp defines the callback interface that a webservice must implement
type Gapp interface {
	// LoadConfig callback allows the app to load the app config. Optionally save the config as a resource for use outside of callbacks
	LoadConfig() gappconfig.Config
	// ConfigureLogging callback allows the app to do any custom log configuration (i.e. based on environmental config)
	ConfigureLogging(conf gappconfig.Config)
	// InitResources callback is where the app would set up DB connections, start internal goroutine daemons, etc.
	InitResources(conf gappconfig.Config)
	// SetHandlers callback allows the app to set the webservice's handlers
	SetHandlers(conf gappconfig.Config) []HandlerMapping
	// SetNotFoundHandler callback allows the app to optionally set a 404 handler
	SetNotFoundHandler(conf gappconfig.Config) http.Handler
	// SetMiddleware callback allows the app to set Negroni middleware handlers. Gapp comes with some handy middleware you can use.
	SetMiddleware(conf gappconfig.Config) []negroni.Handler
	// GetServerConf callback prompts the app for the host and port to listen on. The final return value is the length of time to allow handlers to finish on stop before shutting down the service.
	GetServerConf(conf gappconfig.Config) (host string, port int, gracefulTimeout time.Duration)
	// HandleStart callback is fired right before the service starts listening
	HandleStart(host string, port int)
	// HandleStopped callback is fired after the app has stopped listening. Teardown code should go here.
	HandleStopped()
}

// Run runs a Gapp app object, using its callbacks to configure and fire events. Run blocks until the service is stopped.
func Run(app Gapp) {
	config := app.LoadConfig()
	app.ConfigureLogging(config)
	app.InitResources(config)

	r := mux.NewRouter()

	for _, hm := range app.SetHandlers(config) {
		r.HandleFunc(hm.Route, hm.Handler)
	}
	r.NotFoundHandler = app.SetNotFoundHandler(config)

	n := negroni.New(app.SetMiddleware(config)...)

	n.UseHandler(r)

	host, port, gracefulTimeout := app.GetServerConf(config)
	app.HandleStart(host, port)

	graceful.Run(host+":"+strconv.Itoa(port), gracefulTimeout, n)

	app.HandleStopped()
}
