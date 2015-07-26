package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Hearst-DD/gapp"
	"github.com/Hearst-DD/gappconfig"
	"github.com/codegangsta/negroni"
)

type app struct {
}

func NewExampleApp() *app {
	return &app{}
}

func (a *app) LoadConfig() gappconfig.Config {
	config = gappconfig.New("MYAPP_", gappconfig.Map{
		{"ENV", "dev"},
		{"PRETTY", true},
		{"HOST", "localhost"},
		{"PORT", 4001},
	})
	return config
}

func (a *app) ConfigureLogging(conf gappconfig.Config) {
	// do custom log configuration here...
	logger = log.New(os.Stdout, "MYAPP ", log.LstdFlags)

	logger.Printf("logging configured.")
}

func (a *app) InitResources(conf gappconfig.Config) {
	logger.Printf("initializing...")

	// initialize things like database connections, daemon threads, etc.

	logger.Printf("...done.")
}

func (a *app) SetHandlers(conf gappconfig.Config) []gapp.HandlerMapping {
	return []gapp.HandlerMapping{
		{"/hello/world", helloWorldHandler},
	}
}

func (a *app) SetNotFoundHandler(conf gappconfig.Config) http.Handler {
	// set a not found handler if desired, or use the default
	return nil
}

func (a *app) SetMiddleware(conf gappconfig.Config) []negroni.Handler {
	return []negroni.Handler{
		gapp.RecoveryMiddleware(func(rw http.ResponseWriter, r *http.Request) {
			if err := recover(); err != nil {
				logger.Printf("recovering from panic! err: %v", err)
				http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
			}
		}),
		// only using the post handle logging function
		gapp.LoggingMiddleware(nil, func(method, path string, status int, dur time.Duration) {
			logger.Printf("%s %s completed with %v %s in %v", method, path, status, http.StatusText(status), dur)
		}),
	}
}

func (a *app) GetServerConf(conf gappconfig.Config) (host string, port int, gracefulTimeout time.Duration) {
	host = conf.String("HOST")
	port = conf.Int("PORT")
	gracefulTimeout = time.Second * 10
	return
}

func (a *app) HandleStart(host string, port int) {
	logger.Printf("service started on %s:%d...", host, port)
}

func (a *app) HandleStopped() {
	logger.Printf("...service stopped.")
}
