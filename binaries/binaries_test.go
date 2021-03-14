package binaries

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func tmpDir(t *testing.T) string {
	dir, err := ioutil.TempDir("/tmp", "prisma-client-go-test-fetchEngine-")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestFetch(t *testing.T) {
	dir := tmpDir(t)
	//goland:noinspection GoUnhandledErrorResult
	defer os.RemoveAll(dir)

	if err := FetchNative(dir); err != nil {
		t.Fatalf("fetchEngine failed: %s", err)
	}
}

func TestFetch_localDir(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := FetchNative(wd); err != nil {
		t.Fatalf("fetchEngine failed: %s", err)
	}
}

func TestFetch_withCache(t *testing.T) {
	dir := tmpDir(t)
	//goland:noinspection GoUnhandledErrorResult
	defer os.RemoveAll(dir)

	start := time.Now()
	if err := FetchNative(dir); err != nil {
		t.Fatalf("fetchEngine 1 failed: %s", err)
	}

	log.Printf("first fetchEngine took %s", time.Since(start))

	start = time.Now()
	if err := FetchNative(dir); err != nil {
		t.Fatalf("fetchEngine 2 failed: %s", err)
	}

	log.Printf("second fetchEngine took %s", time.Since(start))

	if time.Since(start) > 20*time.Millisecond {
		t.Fatalf("second fetchEngine took more than 10ms")
	}
}

func TestFetch_relativeDir(t *testing.T) {
	actual := FetchNative(".")
	expected := fmt.Errorf("toDir must be absolute")
	assert.Equal(t, expected, actual)
}
