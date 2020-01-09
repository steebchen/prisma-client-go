package logger

import (
	"io/ioutil"
	"log"
	"os"
)

// TODO add log levels

var v = os.Getenv("PHOTON_GO_LOG")
var Enabled = v != ""

var L *log.Logger

func init() {
	L = log.New(ioutil.Discard, "", 0)
	if Enabled {
		L = log.New(os.Stdout, "", log.Flags())
	}
}
