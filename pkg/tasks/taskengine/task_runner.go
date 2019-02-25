package taskengine

import (
	"fmt"
	"os"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/autoenv/features"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

type TaskRunner interface {
	Run(*taskapi.Context, *taskapi.Task) error
}

type TaskRunnerImpl struct{}

func (r *TaskRunnerImpl) Run(ctx *taskapi.Context, task *taskapi.Task) (err error) {
	for _, action := range task.Actions {
		err = r.runAction(ctx, action)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TaskRunnerImpl) runAction(ctx *taskapi.Context, action taskapi.TaskAction) error {
	desc := action.Description()

	result := action.Needed(ctx)
	if result.Error != nil {
		return fmt.Errorf("The task action (%s) failed to detect whether it need to run: %s", desc, result.Error)
	}

	if result.Needed {
		if desc != "" {
			ctx.UI.TaskActionHeader(desc)
		}
		ctx.UI.Debug("Reason: \"%s\"", result.Reason)

		err := action.Run(ctx)
		if err != nil {
			return fmt.Errorf("The task action failed to run: %s", err)
		}

		result = action.Needed(ctx)
		if result.Error != nil {
			return fmt.Errorf("The task action failed to detect if it is resolved: %s", result.Error)
		}

		if result.Needed {
			return fmt.Errorf("The task action did not produce the expected result: %s", result.Reason)
		}
	}

	feature := action.Feature()
	if feature != nil {
		return r.activateFeature(ctx, *feature)
	}
	return nil
}

func (r *TaskRunnerImpl) activateFeature(ctx *taskapi.Context, feature autoenv.FeatureInfo) error {
	def, err := features.GetRegister().Get(feature.Name)
	if err != nil {
		return err
	}

	devUpNeeded, err := def.Activate(feature.Param, ctx.Cfg, ctx.Project, ctx.Env)
	if err != nil {
		return err
	}
	if devUpNeeded {
		ctx.UI.TaskWarning(fmt.Sprintf("Something is wrong, the feature %s could not be activated", feature))
	}

	// Special case, we want the bud process to get PATH updates from features to call the right processes.
	// Like the pip process from the newly activated virtualenv.
	// Explanation: exec.Command calls exec.LookPath to find the executable path, which rely on the PATH of
	// the process itself.
	// There is no problem when executing a shell command since the shell process will do the executable lookup
	// itself with the PATH value from the specified Env.
	return os.Setenv("PATH", ctx.Env.Get("PATH"))
}
