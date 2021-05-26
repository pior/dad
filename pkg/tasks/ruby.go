package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

const rubyTaskName = "ruby"

func init() {
	taskapi.Register(rubyTaskName, "ruby", parserRuby)
}

func parserRuby(config *taskapi.TaskConfig, task *taskapi.Task) error {
	version, err := config.GetStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	engine := "ruby"
	if config.IsHash() {
		engine, err = config.GetStringPropertyDefault("engine", "ruby")
		if err != nil {
			return err
		}
	}

	task.Info = version

	parserRubyInstallCommand(task)
	parserRubyInstallRubyVersion(task, engine, version)
	return nil
}

func parserRubyInstallCommand(task *taskapi.Task) {
	needed := func(ctx *context.Context) *taskapi.ActionResult {
		result := commandSilent(ctx, "which", "ruby-install").Capture()
		if result.Error != nil {
			return taskapi.ActionNeeded("ruby-install is not installed")
		}
		return taskapi.ActionNotNeeded()
	}
	run := func(ctx *context.Context) error {
		result := command(ctx, "brew", "install", "ruby-install").Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install ruby-install: %s", result.Error)
		}
		return nil
	}
	task.AddActionWithBuilder("install ruby-install", run).OnFunc(needed)
}

func parserRubyInstallRubyVersion(task *taskapi.Task, engine, version string) {
	needed := func(ctx *context.Context) *taskapi.ActionResult {
		rubyInstall := helpers.NewRubyInstall(ctx.Cfg, engine, version)
		if rubyInstall.Installed() {
			return taskapi.ActionNotNeeded()
		}
		return taskapi.ActionNeeded("ruby version is not installed at " + rubyInstall.Path())
	}
	run := func(ctx *context.Context) error {
		rubyInstall := helpers.NewRubyInstall(ctx.Cfg, engine, version)
		return rubyInstall.Install()
	}
	task.AddActionWithBuilder("running ruby-install", run).OnFunc(needed)
}
