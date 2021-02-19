package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCustom(t *testing.T) {
	task := ensureLoadTestTask(t, `
custom:
  met?: test-command
  meet: custom-command
`)

	require.Equal(t, "Task Custom (custom-command) actions=1", task.Describe())
}

func TestCustomName(t *testing.T) {
	task := ensureLoadTestTask(t, `
custom:
  name: NAMENAME
  met?: test-command
  meet: custom-command
`)

	require.Equal(t, "Task Custom (NAMENAME) actions=1", task.Describe())
}

func TestCustomWithBoolean(t *testing.T) {
	_, err := loadTestTask(t, `
custom:
  met?: false
  meet: custom-command
`)

	require.EqualError(t, err, `task "custom": key "met?": expecting a string, found a bool (false)`)
}
