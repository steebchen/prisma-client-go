package logger

import (
	"io/ioutil"
	"log"
	"os"
)

const flag = log.Ldate | log.Lmicroseconds

// TODO add log levels

var v = os.Getenv("PRISMA_CLIENT_GO_LOG")
var Enabled = v != ""

var Debug *log.Logger
var Info *log.Logger

func init() {
	discard := log.New(ioutil.Discard, "", 0)

	Debug = discard
	if Enabled {
		Debug = log.New(os.Stdout, "prisma-client-go debug: ", flag)
	}

	Info = log.New(os.Stdout, "prisma-client-go info: ", flag)
}
