package taskengine

import (
	"fmt"
	"os"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/autoenv/features"
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

type TaskRunner interface {
	Run(*taskapi.Task) error
}

type TaskRunnerImpl struct {
	ctx *context.Context
}

func NewTaskRunner(ctx *context.Context) TaskRunner {
	return &TaskRunnerImpl{ctx: ctx}
}

func (r *TaskRunnerImpl) Run(task *taskapi.Task) (err error) {
	for _, action := range task.Actions {
		err = r.runAction(action)
		if err != nil {
			return fmt.Errorf(`action "%s": %w`, action.Description(), err)
		}
	}
	return nil
}

func (r *TaskRunnerImpl) runAction(action taskapi.TaskAction) error {
	desc := action.Description()

	result := action.Needed(r.ctx)
	if result.Error != nil {
		return fmt.Errorf("detecting whether it needs to run: %w", result.Error)
	}

	if result.Needed {
		if desc != "" {
			r.ctx.UI.TaskActionHeader(desc)
		}
		r.ctx.UI.Debug("Reason: \"%s\"", result.Reason)

		err := action.Run(r.ctx)
		if err != nil {
			return fmt.Errorf("failed to run: %s", err)
		}

		result = action.Needed(r.ctx)
		if result.Error != nil {
			return fmt.Errorf("detecting whether it still needs to run: %s", result.Error)
		}

		if result.Needed {
			return fmt.Errorf("ran successfully but still need to run: %s", result.Reason)
		}
	}

	feature := action.Feature()
	if feature != nil {
		return r.activateFeature(*feature)
	}
	return nil
}

func (r *TaskRunnerImpl) activateFeature(feature autoenv.FeatureInfo) error {
	def, err := features.GlobalRegister().Get(feature.Name)
	if err != nil {
		return err
	}

	devUpNeeded, err := def.Activate(r.ctx, feature.Param)
	if err != nil {
		return err
	}
	if devUpNeeded {
		r.ctx.UI.TaskWarning(fmt.Sprintf("Something is wrong, the feature %s could not be activated", feature))
	}

	// Special case, we want the bud process to get PATH updates from features to call the right processes.
	// Like the pip process from the newly activated virtualenv.
	// Explanation: exec.Command calls exec.LookPath to find the executable path, which rely on the PATH of
	// the process itself.
	// There is no problem when executing a shell command since the shell process will do the executable lookup
	// itself with the PATH value from the specified Env.
	return os.Setenv("PATH", r.ctx.Env.Get("PATH"))
}
