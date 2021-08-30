package logger

import (
	"io/ioutil"
	"log"
	"os"
)

// TODO add log levels

var v = os.Getenv("PHOTON_GO_LOG")
var Enabled = v != ""

var Debug *log.Logger
var Info *log.Logger

func init() {
	discard := log.New(ioutil.Discard, "", 0)

	Debug = discard
	if Enabled {
		Debug = log.New(os.Stdout, "prisma-client-go debug: ", log.Flags())
	}

	Info = log.New(os.Stdout, "prisma-client-go info: ", log.Flags())
}
