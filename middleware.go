package gapp

import (
	"net/http"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/urfave/negroni"
)

type recoveryMiddleware struct {
	recoverFunc func(rw http.ResponseWriter, r *http.Request)
}

type loggingMiddleware struct {
	preLogFunc  func(method, path string, start time.Time)
	postLogFunc func(method, path string, status int, dur time.Duration)
}

type gzipMiddleware struct{}

// RecoveryMiddleware creates a middleware to handle panics during requests with the supplied func
func RecoveryMiddleware(recoverFunc func(rw http.ResponseWriter, r *http.Request)) negroni.Handler {
	return &recoveryMiddleware{
		recoverFunc: recoverFunc,
	}
}

// LoggingMiddleware creates a middleware to log before and after requests. Nil pre or post funcs are OK.
func LoggingMiddleware(preLogFunc func(method, path string, start time.Time),
	postLogFunc func(method, path string, status int, dur time.Duration)) negroni.Handler {

	return &loggingMiddleware{
		preLogFunc:  preLogFunc,
		postLogFunc: postLogFunc,
	}
}

func GzipMiddleware() negroni.Handler {
	return &gzipMiddleware{}
}

func (rec *recoveryMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if rec.recoverFunc != nil {
		defer rec.recoverFunc(rw, r)
	}

	next(rw, r)
}

func (l *loggingMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	if l.preLogFunc != nil {
		l.preLogFunc(r.Method, r.URL.Path, start)
	}

	next(rw, r)

	if l.postLogFunc != nil {
		res := rw.(negroni.ResponseWriter)
		l.postLogFunc(r.Method, r.URL.Path, res.Status(), time.Since(start))
	}
}

func (gm *gzipMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	zippedHandler := gziphandler.GzipHandler(next)
	zippedHandler.ServeHTTP(rw, r)
}
