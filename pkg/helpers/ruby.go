package helpers

import (
	"fmt"
	"os"
	"path"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

type RubyInstall struct {
	engine  string
	version string
	path    string
}

func NewRubyInstall(cfg *config.Config, engine, version string) *RubyInstall {
	return &RubyInstall{
		engine:  engine,
		version: version,
		path:    cfg.DataDir("ruby", engine+"-"+version),
	}
}

func (r *RubyInstall) Installed() bool {
	return utils.PathExists(r.Which("ruby"))
}

func (r *RubyInstall) Path() string {
	return r.path
}

func (r *RubyInstall) Which(program string) string {
	return path.Join(r.path, "bin", program)
}

func (r *RubyInstall) Install() error {
	err := os.MkdirAll(r.path, 0750)
	if err != nil {
		return err
	}

	result := executor.New("ruby-install", "--latest", "--src-dir", "/tmp", "--install-dir", r.path, r.engine, r.version).Run()
	if result.Error != nil {
		return fmt.Errorf("running ruby-install: %w", result.Error)
	}

	return nil
}
