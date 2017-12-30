package hook

import (
	"fmt"
	"os"

	"github.com/pior/dad/pkg/features"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

func Hook() {
	// In the shell hook, the stdout is evaluated by the shell
	// stderr is used to display messages to the user
	ui := termui.NewHookUI()

	proj, err := project.FindCurrent()

	if err != nil && err != project.ErrProjectNotFound {
		// We can't annoy the user here, just quit silently
		ui.Debug("error while searching for current project: %s", err)
		return
	}

	handleFeatures(proj, ui)
}

func handleFeatures(proj *project.Project, ui *termui.UI) {
	env := features.NewEnv(os.Environ())

	var err error
	wantedFeatures := map[string]string{}

	if proj != nil {
		wantedFeatures, err = proj.GetFeatures()
		if err != nil {
			ui.Debug("failed to get project tasks: %s", err)
		}
	}

	activeFeatures := env.GetActiveFeatures()

	for name, featureBuilder := range features.FeatureMap {
		activeVersion, active := activeFeatures[name]
		wantVersion, want := wantedFeatures[name]
		feature := featureBuilder(wantVersion)

		if want {
			if !active || wantVersion != activeVersion {
				err = feature.Enable(proj, env, ui)
				if err != nil {
					if err == features.DevUpNeeded {
						ui.HookWarning("failed to activate %s(%s). Try running dad up first!", name, wantVersion)
					} else {
						ui.Debug("failed: %s", err)
					}
				} else {
					ui.HookFeatureActivated(name, wantVersion)
				}
			}
		} else {
			if active {
				feature.Disable(proj, env, ui)
				ui.Debug("%s deactivated", name)
			}
		}
	}

	env.SetActiveFeatures(wantedFeatures)

	envChanges := env.Changed()
	for _, change := range envChanges {
		ui.Debug("Env change: %+v", change)

		if change.Deleted {
			fmt.Printf("unset %s", change.Name)
		} else {
			fmt.Printf("export %s=\"%s\"", change.Name, change.Value)
		}
	}
}