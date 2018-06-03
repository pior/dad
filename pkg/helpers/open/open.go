package open

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/pior/dad/pkg/helpers"
	"github.com/pior/dad/pkg/project"
)

// Open a file or URL with the default application, return immediately
func Open(location string) error {
	openCommand := "xdg-open"
	if runtime.GOOS == "darwin" {
		openCommand = "open"
	}

	return exec.Command(openCommand, location).Start()
}

// FindLink returns the url of a link about the project
func FindLink(proj *project.Project, linkName string) (url string, err error) {
	if linkName == "" {
		if len(proj.Manifest.Open) == 1 {
			for _, url = range proj.Manifest.Open {
				return url, nil
			}
		}
		return "", fmt.Errorf("more than one link defined: which one do you want?")
	}

	url = proj.Manifest.Open[linkName]
	if url != "" {
		return
	}

	switch linkName {
	case "github", "gh":
		url, err = helpers.NewGitRepo(proj.Path).BuildGithubProjectURL()
		return
	case "pullrequest", "pr":
		url, err = helpers.NewGitRepo(proj.Path).BuildGithubPullrequestURL()
		return
	default:
		err = fmt.Errorf("no link for '%s'", linkName)
	}

	return
}
