package helpers

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

type Crystal struct {
	version string
	path    string
	tarDir  string
}

func NewCrystal(cfg *config.Config, version string) *Crystal {
	return &Crystal{
		version: version,
		path:    cfg.DataDir("crystal", version),
		tarDir:  cfg.DataDir("crystal"),
	}
}

func (c *Crystal) Exists() bool {
	return utils.PathExists(c.Which("crystal"))
}

func (c *Crystal) Path() string {
	return c.path
}

func (c *Crystal) BinPath() string {
	return path.Join(c.path, "bin")
}

func (c *Crystal) Which(program string) string {
	return path.Join(c.path, "bin", program)
}

func (c *Crystal) Archicture() string {
	if runtime.GOARCH == "amd64" {
		return "x86_64"
	}
	return runtime.GOARCH
}

func (c *Crystal) Install() (err error) {
	archiveName := fmt.Sprintf("crystal-%s-1-%s-%s.tar.gz", c.version, runtime.GOOS, c.Archicture())
	tarPath := path.Join(c.tarDir, archiveName)

	if !utils.PathExists(tarPath) {
		err = os.MkdirAll(c.tarDir, 0750)
		if err != nil {
			return
		}

		url := fmt.Sprintf("https://github.com/crystal-lang/crystal/releases/download/%s/%s", c.version, archiveName)
		err = NewDownloader(url).DownloadToFile(tarPath)
		if err != nil {
			return fmt.Errorf("failed to download Crystal %s from %s: %s", c.version, url, err)
		}
	}

	err = os.MkdirAll(c.path, 0750)
	if err != nil {
		return
	}

	result := executor.New("tar", "--strip", "1", "-xzC", c.path, "-f", tarPath).Run()

	if result.Error != nil {
		return fmt.Errorf("failed to extract %s to %s: %s", tarPath, c.path, result.Error)
	}

	return nil
}
