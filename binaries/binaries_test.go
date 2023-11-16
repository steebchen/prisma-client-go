package binaries

import (
	"fmt"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

func tmpDir(t *testing.T) string {
	dir, err := os.MkdirTemp("/tmp", "prisma-client-go-test-fetchEngine-")
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

	dir := path.Join(wd, "out")

	//goland:noinspection GoUnhandledErrorResult
	defer os.RemoveAll(dir)

	if err := FetchNative(dir); err != nil {
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
	massert.Equal(t, expected, actual)
}
