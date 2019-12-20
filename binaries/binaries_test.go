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
	dir, err := ioutil.TempDir("/tmp", "photongo-test-fetch-")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestFetch(t *testing.T) {
	dir := tmpDir(t)
	defer os.RemoveAll(dir)

	if err := FetchNative(dir); err != nil {
		t.Fatalf("fetch failed: %s", err)
	}
}

func TestFetch_localDir(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := FetchNative(wd); err != nil {
		t.Fatalf("fetch failed: %s", err)
	}
}

func TestFetch_withCache(t *testing.T) {
	dir := tmpDir(t)
	defer os.RemoveAll(dir)

	start := time.Now()
	if err := FetchNative(dir); err != nil {
		t.Fatalf("fetch 1 failed: %s", err)
	}

	log.Printf("first fetch took %s", time.Since(start))

	start = time.Now()
	if err := FetchNative(dir); err != nil {
		t.Fatalf("fetch 2 failed: %s", err)
	}

	log.Printf("second fetch took %s", time.Since(start))

	if time.Since(start) > 10*time.Millisecond {
		t.Fatalf("second fetch took more than 10ms")
	}
}

func TestFetch_relativeDir(t *testing.T) {
	err := FetchNative(".")
	expected := fmt.Errorf("toDir must be absolute")
	assert.Equal(t, expected, err)
}
