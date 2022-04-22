package download

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {

	// Create dummy http server
	server := makeDelayedServer(0 * time.Millisecond)
	defer server.Close()

	// Create dummy download folder
	dir, err := os.MkdirTemp("", "TestAddUrl")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	t.Run("Distinct filenames", func(t *testing.T) {

		d := New()

		if err := d.Add(server.URL, "_unique1_"); err != nil {
			t.Error(err)
		}
		if err := d.Add(server.URL, "_unique2_"); err != nil {
			t.Error(err)
		}
		if err := d.Add(server.URL, "_unique3_"); err != nil {
			t.Error(err)
		}
	})

	t.Run("Duplicate filenames", func(t *testing.T) {

		d := New()

		if err := d.Add(server.URL, "_dup1_"); err != nil {
			t.Error(err)
		}
		if err := d.Add(server.URL, "_dup1_"); !errors.Is(err, errDuplicateFilename) {
			t.Error(err)
		}
	})
}

func TestRun(t *testing.T) {

	// Create dummy http server
	server := makeDelayedServer(0 * time.Millisecond)
	defer server.Close()

	// Create dummy download folder
	dir, err := os.MkdirTemp("", "TestRun")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	d := New()
	err = d.DestDir(dir)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10; i++ {
		d.Add(server.URL, fmt.Sprintf("_testrun_%v", i))
	}
	d.Run()

	// Check the dir/_testrun_ exists and its contents are httpBody
}

var httpBody = "Hello"

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(httpBody))
	}))
}
