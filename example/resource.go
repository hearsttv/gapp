package main

import (
	"log"

	"github.com/Hearst-DD/gappconfig"
)

var logger *log.Logger

// save a reference to the loaded config if needed (e.g. in a handler function)
var config gappconfig.Config
