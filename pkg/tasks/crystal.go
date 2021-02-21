package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.Register("crystal", "Crystal", parseCrystal)
}

func parseCrystal(config *taskapi.TaskConfig, task *taskapi.Task) error {
	version, err := config.GetStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	task.Info = version

	installNeeded := func(ctx *context.Context) *taskapi.ActionResult {
		if !helpers.NewCrystal(ctx.Cfg, version).Exists() {
			return taskapi.ActionNeeded("crystal distribution is not installed")
		}
		return taskapi.ActionNotNeeded()
	}
	installCrystal := func(ctx *context.Context) error {
		return helpers.NewCrystal(ctx.Cfg, version).Install()
	}
	task.AddActionWithBuilder("install crystal distribution", installCrystal).
		OnFunc(installNeeded).
		SetFeature("crystal", version)

	return nil
}
