package utils

import (
	_ "embed"
	"runtime/debug"
)

var Commit = func() string {
	var commitMessage string = "unknown"
	var modified bool = false

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commitMessage = setting.Value
			} else if setting.Key == "vcs.modified" {
				modified = setting.Value == "true"
			}
		}
	}

	if modified {
		commitMessage = commitMessage + "*"
	}

	return commitMessage
}()
