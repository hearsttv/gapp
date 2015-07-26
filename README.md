## Gapp: A simple webservice framework for Go services

Gapp allows you to quickly author webservices without writing server boilerplate each time. It was developed for internal use by the Hearst Digital Product Development Group to streamline the creation and maintanence of Go services.

## Install

```bash
go get github.com/Hearst-DD/gapp
```

NOTE: Gapp is used interally by Hearst DPDG and may change. If you want to use Gapp for your service, it is recommended that you vendor or fork the code.

## Dependencies

Gapp uses the following popular tools to construct webservices: 
* Gorilla Mux (http://www.gorillatoolkit.org/pkg/mux)
* Negroni middleware (https://github.com/codegangsta/negroni)
* Graceful (https://github.com/tylerb/graceful) 

## Use

Gapp works by providing a system of callbacks to initialize, run and manage a webservice. To use Gapp to power your service, create an app object that implements the interface. 

Gapp callbacks suggest (but do not enforce) a clear file structure for your app: 
* main.go simply instantiates and runs your app
* app.go holds your implementation of the Gapp interface
* resource.go contains package level vars and functions for resources such as loggers, database connections, etc.
* handler.go contains your handler definitions
* subpackages contain internal service functionality (DB wrappers, daemon goroutines, etc.)

See the following example for details.

## Example

The example/ directory contains a full implementation of a Gapp app. To run it, cd into the example directory, build, run and visit localhost:4001/hello/world in your browser.

Sample output: 
```bash
Nathans-MacBook-Pro:example nate$ go install .
Nathans-MacBook-Pro:example nate$ example 
MYAPP 2015/07/26 12:33:54 logging configured.
MYAPP 2015/07/26 12:33:54 initializing...
MYAPP 2015/07/26 12:33:54 ...done.
MYAPP 2015/07/26 12:33:54 service started on localhost:4001...
MYAPP 2015/07/26 12:33:57 GET /hello/world completed with 200 OK in 23.21µs
^C
MYAPP 2015/07/26 12:34:13 ...service stopped.
```
