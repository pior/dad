package tasks

import (
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	t := registerTaskDefinition("go")
	t.name = "Golang"
	t.parser = parseGolang
}

func parseGolang(config *taskConfig, task *Task) error {
	version, err := config.getStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	task.header = version
	task.featureName = "golang"
	task.featureParam = version

	task.addAction(&golangGoPath{})
	task.addAction(&golangInstall{version: version})

	return nil
}

type golangGoPath struct{}

func (g *golangGoPath) description() string {
	return ""
}

func (g *golangGoPath) needed(ctx *Context) *actionResult {
	if ctx.env.Get("GOPATH") == "" {
		return actionNeeded("GOPATH is not set")
	}

	if !strings.Contains(ctx.env.Get("PATH"), g.binPath(ctx)) {
		return actionNeeded(fmt.Sprintf("%s is not in PATH", g.binPath(ctx)))
	}

	return actionNotNeeded()
}

func (g *golangGoPath) run(ctx *Context) error {
	ctx.ui.TaskWarning("The GOPATH environment variable should be set to ~/")
	return nil
}

func (g *golangGoPath) binPath(ctx *Context) string {
	return fmt.Sprintf("%s/bin", ctx.env.Get("GOPATH"))
}

type golangInstall struct {
	version string
}

func (g *golangInstall) description() string {
	return fmt.Sprintf("Install Go version %s", g.version)
}

func (g *golangInstall) needed(ctx *Context) *actionResult {
	if !helpers.NewGolang(ctx.cfg, g.version).Exists() {
		return actionNeeded("golang distribution is not installed")
	}
	return actionNotNeeded()
}

func (g *golangInstall) run(ctx *Context) error {
	return helpers.NewGolang(ctx.cfg, g.version).Install()
}
