package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/hearsttv/gapp"
)

type app struct {
}

func NewExampleApp() *app {
	return &app{}
}

func (a *app) LoadConfig() gapp.Config {
	config = gapp.NewConfig("MYAPP_", gapp.ConfigMap{
		{"ENV", "dev"},
		{"PRETTY", true},
		{"HOST", "localhost"},
		{"PORT", 4001},
	})
	return config
}

func (a *app) ConfigureLogging(conf gapp.Config) {
	// do custom log configuration here...
	logger = log.New(os.Stdout, "MYAPP ", log.LstdFlags)

	logger.Printf("logging configured.")
}

func (a *app) InitResources(conf gapp.Config) {
	logger.Printf("initializing...")

	// initialize things like database connections, daemon threads, etc.

	logger.Printf("...done.")
}

func (a *app) ConfigureRoutes(r *mux.Router, conf gapp.Config) {
	r.HandleFunc("/hello/world", helloWorldHandler).Methods("GET")

	// set a not found handler if desired, or use the default
}

func (a *app) SetMiddleware(conf gapp.Config) []negroni.Handler {
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

func (a *app) GetServerConf(conf gapp.Config) gapp.ServerConfig {
	return gapp.ServerConfig{
		Host:            conf.String("HOST"),
		Port:            conf.Int("PORT"),
		GracefulTimeout: time.Second * 10,
	}
}

func (a *app) HandleStart(host string, port, tlsPort int) {
	logger.Printf("service started on %s:%d...", host, port)
}

func (a *app) HandleStopped() {
	logger.Printf("...service stopped.")
}
