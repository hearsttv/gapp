package gapp

import (
	"net/http"
	"strconv"
	"sync"
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

// ServerConfig encapsulates the various values needed to start the server
type ServerConfig struct {
	Host              string
	Port              int
	GracefulTimeout   time.Duration
	HTTPSEnabled      bool
	HTTPDisabled      bool
	TLSPort           int
	TLSCertFile       string
	TLSPrivateKeyFile string
}

// Gapp defines the callback interface that a webservice must implement
type Gapp interface {
	// LoadConfig callback allows the app to load the app config. Optionally save the config as a resource for use outside of callbacks
	LoadConfig() gappconfig.Config
	// ConfigureLogging callback allows the app to do any custom log configuration (i.e. based on environmental config)
	ConfigureLogging(conf gappconfig.Config)
	// InitResources callback is where the app would set up DB connections, start internal goroutine daemons, etc.
	InitResources(conf gappconfig.Config)
	// ConfigureRoutes callback allows the app to set the webservice's handlers
	ConfigureRoutes(r *mux.Router, conf gappconfig.Config)
	// SetMiddleware callback allows the app to set Negroni middleware handlers. Gapp comes with some handy middleware you can use.
	SetMiddleware(conf gappconfig.Config) []negroni.Handler
	// GetServerConf callback prompts the app for the host and port to listen on. The final return value is the length of time to allow handlers to finish on stop before shutting down the service.
	GetServerConf(conf gappconfig.Config) ServerConfig
	// HandleStart callback is fired right before the service starts listening
	HandleStart(host string, port int)
	// HandleStopped callback is fired after the app has stopped listening. Teardown code should go here.
	HandleStopped()
}

// Run runs a Gapp app object, using its callbacks to configure and fire events. Run blocks until the service is stopped.
func Run(app Gapp) {
	config, n := initApp(app)

	serverConfig := app.GetServerConf(config)
	app.HandleStart(serverConfig.Host, serverConfig.Port)

	if serverConfig.HTTPDisabled && !serverConfig.HTTPSEnabled {
		panic("must accept at least one scheme (HTTP and/or HTTPS)")
	}

	var wg sync.WaitGroup

	if serverConfig.HTTPSEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()

			srv := &graceful.Server{
				Timeout: serverConfig.GracefulTimeout,

				Server: &http.Server{
					Addr:    serverConfig.Host + ":" + strconv.Itoa(serverConfig.TLSPort),
					Handler: n,
				},
			}
			err := srv.ListenAndServeTLS(serverConfig.TLSCertFile, serverConfig.TLSPrivateKeyFile)
			if err != nil {
				panic(err)
			}
		}()
	}
	if !serverConfig.HTTPDisabled {
		wg.Add(1)
		go func() {
			defer wg.Done()

			srv := &graceful.Server{
				Timeout: serverConfig.GracefulTimeout,

				Server: &http.Server{
					Addr:    serverConfig.Host + ":" + strconv.Itoa(serverConfig.Port),
					Handler: n,
				},
			}
			srv.ListenAndServe()
		}()
	}

	wg.Wait()

	app.HandleStopped()
}

var doRunFunc = graceful.Run

func initApp(app Gapp) (gappconfig.Config, *negroni.Negroni) {
	config := app.LoadConfig()
	app.ConfigureLogging(config)
	app.InitResources(config)

	r := mux.NewRouter()

	app.ConfigureRoutes(r, config)

	n := negroni.New(app.SetMiddleware(config)...)

	n.UseHandler(r)

	return config, n
}
