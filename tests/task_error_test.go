package integration

import (
	"testing"

	"github.com/devbuddy/devbuddy/tests/context"
)

func Test_Task_Error_NotAList(t *testing.T) {
	c := CreateContextAndInit(t)

	c.Write("dev.yml", `up: somestring`)

	lines := c.Run("bud up", context.ExitCode(1))
	OutputEqual(t, lines,
		"Error: yaml: unmarshal errors:",
		"  line 1: cannot unmarshal !!str `somestring` into []interface {}",
	)
}

func Test_Task_Error_InvalidType(t *testing.T) {
	c := CreateContextAndInit(t)

	c.Write("dev.yml", `up: [true]`)

	lines := c.Run("bud up", context.ExitCode(1))
	OutputEqual(t, lines, `Error: invalid task: "true"`)
}

func Test_Task_Error_UnknownTask(t *testing.T) {
	c := CreateContextAndInit(t)

	c.Write("dev.yml", `up: [notatask]`)

	lines := c.Run("bud up", context.ExitCode(0)) // TODO: This should probably return 1
	OutputEqual(t, lines,
		"◼︎ Unknown",
		`  Warning: Unknown task: "notatask"`,
	)
}

func Test_Task_Error_Invalid_Hash_Type(t *testing.T) {
	c := CreateContextAndInit(t)

	c.Write("dev.yml", `up: [ { go: {version: 1.16} } ]`)

	lines := c.Run("bud up", context.ExitCode(1))
	OutputEqual(t, lines,
		`Error: task "go": key "version": expecting a string, not a float64 (1.16)`,
	)
}

func Test_Task_Error_Invalid_List(t *testing.T) {
	c := CreateContextAndInit(t)

	c.Write("dev.yml", `up: [ { homebrew: {} } ]`)

	lines := c.Run("bud up", context.ExitCode(1))
	OutputEqual(t, lines,
		`Error: task "homebrew": expecting a list of strings, not a map[interface {}]interface {} (map[])`,
	)
}
