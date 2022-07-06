// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"github.com/kayuii/github-checker/plugin"

	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(settings *plugin.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "github-url",
			Usage:       "github url, defaults to current scm",
			Aliases:     []string{"u"},
			EnvVars:     []string{"PLUGIN_GITHUB_URL", "DRONE_REPO_LINK"},
			Destination: &settings.GitHubURL,
		},
		&cli.BoolFlag{
			Name:        "PreRelease",
			Usage:       "Pre-Release",
			Aliases:     []string{"p"},
			EnvVars:     []string{"PLUGIN_PRE_RELEASE", "GITHUB_PRE_RELEASE"},
			Destination: &settings.PreRelease,
		},
		&cli.StringFlag{
			Name:        "check-version",
			Usage:       "compare version number",
			Aliases:     []string{"c"},
			EnvVars:     []string{"PLUGIN_VERSION", "GITHUB_RELEASE_VERSION"},
			Destination: &settings.Version,
		},
	}
}
