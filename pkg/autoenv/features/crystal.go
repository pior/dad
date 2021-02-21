package features

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register("crystal", crystalActivate, nil)
}

func crystalActivate(ctx *context.Context, version string) (bool, error) {
	ctx.UI.Debug("Activating crystal")

	crystal := helpers.NewCrystal(ctx.Cfg, version)
	if !crystal.Exists() {
		ctx.UI.Debug("Cannot activate crystal: not installed")
		return true, nil
	}

	ctx.UI.Debug(fmt.Sprintf("Adding %s to the path", crystal.BinPath()))
	ctx.Env.PrependToPath(crystal.BinPath())
	return false, nil
}
