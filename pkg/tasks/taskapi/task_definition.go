package taskapi

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/context"
)

type taskParser func(*TaskConfig, *Task) error

type TaskDefinition struct {
	Key           string     // the value used in dev.yml to identify a task
	Name          string     // the displayed name of the task
	RequiredTask  string     // another task that should be declared before this task
	Parser        taskParser // the config parser for this task
	OSRequirement string     // The platform this task can run on. "debian", "macos"
}

var taskDefinitions = make(map[string]*TaskDefinition)

func Register(key, name string, parserFunc taskParser) *TaskDefinition {
	if _, ok := taskDefinitions[key]; ok {
		panic(fmt.Sprint("Can't re-register a taskDefinition:", name))
	}
	if key == "" || name == "" {
		panic("key and name cannot be empty")
	}

	td := &TaskDefinition{Key: key, Name: name, Parser: parserFunc}
	taskDefinitions[key] = td
	return td
}

func (t *TaskDefinition) SetRequiredTask(name string) *TaskDefinition {
	if name == "" {
		panic("name cannot be empty")
	}
	t.RequiredTask = name
	return t
}

func (t *TaskDefinition) SetOSRequirement(requirement string) *TaskDefinition {
	if requirement == "" {
		panic("requirement cannot be empty")
	}
	t.OSRequirement = requirement
	return t
}

func GetDefinitionOrUnknown(name string) *TaskDefinition {
	taskDef := taskDefinitions[name]
	if taskDef == nil {
		taskDef = newUnknownTaskDefinition()
	}
	return taskDef
}

func newUnknownTaskDefinition() *TaskDefinition {
	parser := func(config *TaskConfig, task *Task) error {
		task.AddActionWithBuilder("", func(ctx *context.Context) error {
			ctx.UI.TaskWarning(fmt.Sprintf("Unknown task: \"%s\"", config.name))
			return nil
		})
		return nil
	}
	return &TaskDefinition{Name: "Unknown", Parser: parser}
}
