package gapp

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/stretchr/testify/assert"
)

func Test_LoggingMiddleware(t *testing.T) {
	state := "start"
	t0 := time.Now()
	var t1 time.Time

	prelogFunc := func(method, path string, start time.Time) {
		assert.Equal(t, "start", state)
		state = "prelogCalled"
		assert.Equal(t, "GET", method)
		assert.Equal(t, "/foo/bar", path)
		assert.True(t, t0.Before(start))
		t1 = start
	}

	next := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "prelogCalled", state)
		state = "nextCalled"
		w.WriteHeader(200)
	}

	postlogFunc := func(method, path string, status int, dur time.Duration) {
		assert.Equal(t, "nextCalled", state)
		state = "postlogCalled"
		assert.Equal(t, "GET", method)
		assert.Equal(t, "/foo/bar", path)
		assert.Equal(t, 200, status)
		assert.True(t, dur <= time.Since(t1))
	}

	r, err := http.NewRequest("GET", "http://example.com/foo/bar", nil)
	assert.Nil(t, err) // sanity
	rw := negroni.NewResponseWriter(httptest.NewRecorder())

	middleware := LoggingMiddleware(prelogFunc, postlogFunc)

	assert.Equal(t, reflect.ValueOf(prelogFunc), reflect.ValueOf(middleware.(*loggingMiddleware).preLogFunc))
	assert.Equal(t, reflect.ValueOf(postlogFunc), reflect.ValueOf(middleware.(*loggingMiddleware).postLogFunc))

	middleware.ServeHTTP(rw, r, next)

	assert.Equal(t, "postlogCalled", state)

}
