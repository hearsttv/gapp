package gapp

import (
	"net/http"
	"testing"
	"time"

	"github.com/Hearst-DD/gappconfig"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var conf = gappconfig.New("FOO_", gappconfig.Map{
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

func (a *testapp) LoadConfig() gappconfig.Config {
	assert.Equal(a.t, "start", state)
	state = "LoadConfigCalled"
	return conf
}

func (a *testapp) ConfigureLogging(conf gappconfig.Config) {
	assert.Equal(a.t, "LoadConfigCalled", state)
	state = "ConfigureLoggingCalled"
}

func (a *testapp) InitResources(conf gappconfig.Config) {
	assert.Equal(a.t, "ConfigureLoggingCalled", state)
	state = "InitResourcesCalled"
}

func (a *testapp) ConfigureRoutes(r *mux.Router, conf gappconfig.Config) {
	assert.Equal(a.t, "InitResourcesCalled", state)
	state = "ConfigureRoutesCalled"
}

func (a *testapp) SetMiddleware(conf gappconfig.Config) []negroni.Handler {
	assert.Equal(a.t, "ConfigureRoutesCalled", state)
	state = "SetMiddlewareCalled"
	return handlers
}

func (a *testapp) GetServerConf(conf gappconfig.Config) (host string, port int, gracefulTimeout time.Duration) {
	return "localhost", 8080, time.Minute
}

func (a *testapp) HandleStart(host string, port int) {
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
