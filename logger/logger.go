package logger

import (
	"io/ioutil"
	"log"
	"os"
)

// TODO move logger from gotpl to here

var Debug = os.Getenv("PHOTON_GO_LOG") != ""

var L *log.Logger

func init() {
	L = log.New(ioutil.Discard, "", 0)
	if Debug {
		L = log.New(os.Stdout, "", log.Flags())
	}
}
