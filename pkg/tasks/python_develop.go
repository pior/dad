package tasks

import (
	"fmt"
	"strings"
)

func init() {
	t := registerTaskDefinition("python_develop")
	t.name = "Python develop"
	t.requiredTask = pythonTaskName
	t.parser = parserPythonDevelop
}

func parserPythonDevelop(config *taskConfig, task *Task) error {
	extras, err := config.getListOfStringsPropertyDefault("extras", []string{})
	if err != nil {
		return err
	}

	pipTarget := "."
	if len(extras) > 0 {
		pipTarget = fmt.Sprintf(".[%s]", strings.Join(extras, ","))
	}
	pipArgs := []string{"install", "--require-virtualenv", "-e", pipTarget}

	builder := actionBuilder("install python package in develop mode", func(ctx *Context) error {
		result := command(ctx, "pip", pipArgs...).AddOutputFilter("already satisfied").Run()
		if result.Error != nil {
			return fmt.Errorf("Pip failed: %s", result.Error)
		}

		return nil
	})
	builder.OnFileChange("setup.py")
	builder.OnTaskChange(pythonTaskName)

	task.addAction(builder.Build())
	return nil
}
