package taskapi

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/autoenv"
	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/project"
)

func GetTasksFromProject(proj *project.Project) (taskList []*Task, err error) {
	var task *Task

	manifest, err := manifest.Load(proj.Path)
	if err != nil {
		return nil, err
	}

	if len(manifest.Env) != 0 {
		task, err = NewTaskFromPayload("env")
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, task)
	}

	for _, payload := range manifest.Up {
		task, err = NewTaskFromPayload(payload)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, task)
	}

	return taskList, nil
}

func NewTaskFromPayload(payload interface{}) (*Task, error) {
	taskConfig, err := NewTaskConfig(payload)
	if err != nil {
		return nil, fmt.Errorf("parsing task: %w", err)
	}

	taskDef := GetDefinitionOrUnknown(taskConfig.name)

	task := &Task{TaskDefinition: taskDef}

	err = taskDef.Parser(taskConfig, task)
	if err != nil {
		return nil, fmt.Errorf(`task "%s": %w`, task.Key, err)
	}

	return task, nil
}

func GetFeaturesFromTasks(tasks []*Task) autoenv.FeatureSet {
	featureSet := autoenv.NewFeatureSet()

	for _, task := range tasks {
		for _, action := range task.Actions {
			if action.Feature() != nil {
				featureSet = featureSet.With(action.Feature())
			}
		}
	}

	return featureSet
}
