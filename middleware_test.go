package gapp

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
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

func Test_LoggingMiddleware_nilFuncs(t *testing.T) {
	state := "start"

	next := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "start", state)
		state = "nextCalled"
	}

	r, err := http.NewRequest("GET", "http://example.com/foo/bar", nil)
	assert.Nil(t, err) // sanity
	rw := negroni.NewResponseWriter(httptest.NewRecorder())

	middleware := LoggingMiddleware(nil, nil)

	assert.Nil(t, middleware.(*loggingMiddleware).preLogFunc)
	assert.Nil(t, middleware.(*loggingMiddleware).postLogFunc)

	middleware.ServeHTTP(rw, r, next)

	assert.Equal(t, "nextCalled", state)
}

func Test_RecoveryMiddleware(t *testing.T) {
	called := false

	recoverFunc := func(rw http.ResponseWriter, r *http.Request) {
		assert.False(t, called)
		called = true
	}

	next := func(w http.ResponseWriter, r *http.Request) {
		// recover func is defered until after next()
		assert.False(t, called)
	}

	r, err := http.NewRequest("GET", "http://example.com/foo/bar", nil)
	assert.Nil(t, err) // sanity
	rw := negroni.NewResponseWriter(httptest.NewRecorder())

	middleware := RecoveryMiddleware(recoverFunc)

	assert.Equal(t, reflect.ValueOf(recoverFunc), reflect.ValueOf(middleware.(*recoveryMiddleware).recoverFunc))

	middleware.ServeHTTP(rw, r, next)

	assert.True(t, called)
}

func Test_RecoveryMiddleware_nilRecoverFunc(t *testing.T) {
	nextCalled := false
	next := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}

	r, err := http.NewRequest("GET", "http://example.com/foo/bar", nil)
	assert.Nil(t, err) // sanity
	rw := negroni.NewResponseWriter(httptest.NewRecorder())

	middleware := RecoveryMiddleware(nil)

	assert.Nil(t, middleware.(*recoveryMiddleware).recoverFunc)

	middleware.ServeHTTP(rw, r, next)

	assert.True(t, nextCalled)
}
