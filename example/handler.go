package main

import (
	"fmt"
	"net/http"
)

func helloWorldHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "Hello, World!")
}
