package tasks

import (
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/helpers/store"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

type genericTaskAction struct {
	builder   *genericTaskActionBuilder
	runCalled bool
}

func (a *genericTaskAction) description() string {
	return a.builder.desc
}

func (a *genericTaskAction) needed(ctx *Context) (result *actionResult) {
	if a.runCalled {
		return a.post(ctx)
	}
	return a.pre(ctx)
}

func (a *genericTaskAction) run(ctx *Context) error {
	a.runCalled = true
	return a.builder.runFunc(ctx)
}

// internals

func (a *genericTaskAction) pre(ctx *Context) (result *actionResult) {
	hasConditions := false

	for _, condition := range a.builder.conditions {
		hasConditions = true

		result = condition.pre(ctx)
		if result.Error != nil || result.Needed {
			return result
		}
	}

	for _, filePath := range a.builder.monitoredFiles {
		hasConditions = true

		result = genericTaskActionPreConditionForFile(ctx, filePath)
		if result.Error != nil || result.Needed {
			return result
		}
	}

	if hasConditions {
		return actionNotNeeded()
	}
	return actionNeeded("action without conditions")
}

func (a *genericTaskAction) post(ctx *Context) (result *actionResult) {
	for _, condition := range a.builder.conditions {
		result = condition.post(ctx)
		if result.Error != nil || result.Needed {
			return result
		}
	}
	for _, filePath := range a.builder.monitoredFiles {
		result = genericTaskActionPostConditionForFile(ctx, filePath)
		if result.Error != nil || result.Needed {
			return result
		}
	}
	return actionNotNeeded()
}

func genericTaskActionPreConditionForFile(ctx *Context, path string) *actionResult {
	fullPath := filepath.Join(ctx.proj.Path, path)

	if !utils.PathExists(fullPath) {
		return actionNeeded("file %s does not exist", path)
	}

	fileChecksum, err := utils.FileChecksum(fullPath)
	if err != nil {
		return actionFailed("failed to get the file checksum: %s", err)
	}

	storedChecksum, err := store.New(ctx.proj.Path).GetString("checksum" + path)
	if err != nil {
		return actionFailed("failed to read the previous file checksum: %s", err)
	}

	if fileChecksum != storedChecksum {
		return actionNeeded("file %s has changed", path)
	}
	return actionNotNeeded()
}

func genericTaskActionPostConditionForFile(ctx *Context, path string) *actionResult {
	fullPath := filepath.Join(ctx.proj.Path, path)

	if !utils.PathExists(fullPath) {
		return actionNotNeeded()
	}

	fileChecksum, err := utils.FileChecksum(fullPath)
	if err != nil {
		return actionFailed("failed to get the file checksum: %s", err)
	}

	err = store.New(ctx.proj.Path).SetString("checksum"+path, fileChecksum)
	if err != nil {
		return actionFailed("failed to store the current file checksum: %s", err)
	}

	return actionNotNeeded()
}

func genericTaskActionPreConditionForTask(ctx *Context, path string) *actionResult {
	storedChecksum, err := store.New(ctx.proj.Path).GetString("checksum" + path)
	if err != nil {
		return actionFailed("failed to read the previous file checksum: %s", err)
	}

	if fileChecksum != storedChecksum {
		return actionNeeded("file %s has changed", path)
	}
	return actionNotNeeded()
}

func genericTaskActionPostConditionForTask(ctx *Context, path string) *actionResult {
}
