package helpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageIsInCellar(t *testing.T) {
	prefix, err := ioutil.TempDir("/tmp", "dad-brew")
	require.NoError(t, err, "ioutil.TempDir() failed")

	cellarPath := filepath.Join(prefix, "Cellar")

	os.MkdirAll(filepath.Join(cellarPath, "curl", "1.2.3"), os.ModePerm)

	h := NewHomebrewWithPrefix(prefix)

	require.Truef(t, h.PackageIsInCellar("curl"), "Curl is missing from Cellar %s", cellarPath)
	require.Falsef(t, h.PackageIsInCellar("vim"), "Curl is missing from Cellar %s", cellarPath)
}
