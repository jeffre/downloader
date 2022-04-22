package download

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
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

	fileCount := 10

	for i := 0; i < fileCount; i++ {
		d.Add(server.URL, fmt.Sprintf("_testrun_%v", i))
	}
	d.Run()

	// Check the dir/_testrun_ exists and its contents are httpBody
	for i := 0; i < fileCount; i++ {
		path := path.Join(dir, fmt.Sprintf("_testrun_%v", i))
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			t.Error(err)
		}
		got := string(bytes)

		if got != httpBody {
			t.Errorf("got %q wanted %q\n", got, httpBody)
		}
	}
}

var httpBody = "Hello"

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(httpBody))
	}))
}
