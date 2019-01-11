package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	t := registerTaskDefinition("node")
	t.name = "NodeJS"
	t.parser = parseNode
}

func parseNode(config *TaskConfig, task *Task) error {
	version, err := config.getStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	task.header = version
	task.feature = features.NewFeatureInfo("node", version)

	run := func(ctx *Context) error {
		return helpers.NewNode(ctx.cfg, version).Install()
	}
	condition := func(ctx *Context) *actionResult {
		if !helpers.NewNode(ctx.cfg, version).Exists() {
			return actionNeeded("node version is not installed")
		}
		return actionNotNeeded()
	}
	task.AddActionWithBuilder("install nodejs from https://nodejs.org", run).OnFunc(condition)

	return nil
}
