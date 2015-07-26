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
	LoadConfig() gappconfig.Config
	ConfigureLogging(conf gappconfig.Config)
	InitResources(conf gappconfig.Config)
	SetHandlers(conf gappconfig.Config) []HandlerMapping
	SetNotFoundHandler(conf gappconfig.Config) http.Handler
	SetMiddleware(conf gappconfig.Config) []negroni.Handler
	GetServerConf(conf gappconfig.Config) (host string, port int, gracefulTimeout time.Duration)
	HandleStart(host string, port int)
	HandleStopped()
}

// Run runs a Gapp app object, using its callbacks to configure and fire events
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
