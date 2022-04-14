package download_test

import (
	"testing"

	"github.com/jeffre/download"
)

func TestAddUrl(t *testing.T) {
	t.Run("1 URL", func(t *testing.T) {
		d := download.New()
		d.Urls = []string{"a"}
		d.Start()
	})

	t.Run("2 URLs", func(t *testing.T) {
		d := download.New()
		d.Urls = []string{"a", "b"}
		d.Start()
	})

	t.Run("multiple URLs", func(t *testing.T) {
		d := download.New()
		d.Urls = []string{"a", "b", "c", "d", "e", "f", "g"}
		d.Start()
	})
}
