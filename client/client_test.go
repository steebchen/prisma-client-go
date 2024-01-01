package client

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	// test prefetch case
	prefetch := []string{"prefetch"}
	assert.Nil(t, Process(prefetch))

	// test init case and remove prisma directory
	init := []string{"init"}
	assert.Nil(t, Process(init))
	pwd, _ := os.Getwd()
	if assert.DirExists(t, pwd+"/prisma") {
		os.Remove(pwd + "/.env")
		os.RemoveAll(pwd + "/prisma")
	}

	// test no args
	empty := []string{}
	assert.Nil(t, Process(empty))

	// test wrong input
	assert.Error(t, Process([]string{"wrong"}))

}
