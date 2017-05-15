package gapp

import (
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

var conf = NewConfig("FOO_", ConfigMap{
	{"BAR", "baz"},
})

var handlers = []negroni.Handler{
	RecoveryMiddleware(func(rw http.ResponseWriter, r *http.Request) {}),
}

// used for asserting order of API methods called
var state = "start"

type testapp struct {
	t *testing.T
}

func (a *testapp) LoadConfig() Config {
	assert.Equal(a.t, "start", state)
	state = "LoadConfigCalled"
	return conf
}

func (a *testapp) ConfigureLogging(conf Config) {
	assert.Equal(a.t, "LoadConfigCalled", state)
	state = "ConfigureLoggingCalled"
}

func (a *testapp) InitResources(conf Config) {
	assert.Equal(a.t, "ConfigureLoggingCalled", state)
	state = "InitResourcesCalled"
}

func (a *testapp) ConfigureRoutes(r *mux.Router, conf Config) {
	assert.Equal(a.t, "InitResourcesCalled", state)
	state = "ConfigureRoutesCalled"
}

func (a *testapp) SetMiddleware(conf Config) []negroni.Handler {
	assert.Equal(a.t, "ConfigureRoutesCalled", state)
	state = "SetMiddlewareCalled"
	return handlers
}

func (a *testapp) GetServerConf(conf Config) ServerConfig {
	return ServerConfig{
		Host:            "localhost",
		Port:            8080,
		GracefulTimeout: time.Minute,
	}
}

func (a *testapp) HandleStart(host string, port, tlsPort int) {
	// noop
}

func (a *testapp) HandleStopped() {
	// noop
}

func Test_initApp(t *testing.T) {
	app := &testapp{t: t}
	initApp(app)

	assert.Equal(t, "SetMiddlewareCalled", state)
}
