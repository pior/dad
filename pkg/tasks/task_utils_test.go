package tasks

import (
	"testing"

	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)

func loadTestTask(t *testing.T, payload string) (*taskapi.Task, error) {
	var data interface{}
	err := yaml.Unmarshal([]byte(payload), &data)
	require.NoError(t, err, "Failed to load a yaml task fixture: %s")
	return taskapi.NewTaskFromPayload(data)
}

func ensureLoadTestTask(t *testing.T, payload string) *taskapi.Task {
	task, err := loadTestTask(t, payload)
	require.NoError(t, err, "NewTaskFromPayload() failed")
	return task
}
