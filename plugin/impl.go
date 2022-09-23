// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/blang/semver/v4"
)

var BaseURL = "https://api.github.com/repos/%s/releases"

// Settings for the plugin.
type Settings struct {
	GitHubURL  string
	Version    string
	Pipe       string
	PreRelease bool
	Mode       int

	baseURL *url.URL
	version semver.Version
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	var err error

	p.settings.baseURL, err = gitHubURLs(p.settings.GitHubURL)
	if err != nil {
		return fmt.Errorf("failed to get GitHub urls: %w", err)
	}

	if len(p.settings.Version) > 0 {
		p.settings.version = Version(p.settings.Version)
		err := p.settings.version.Validate()
		if err != nil {
			return fmt.Errorf("failed to parse version: %w", err)
		}
	} else if len(p.settings.Pipe) > 0 {
		p.settings.version = Version(p.settings.Pipe)
		err := p.settings.version.Validate()
		if err != nil {
			return fmt.Errorf("failed to parse version from pipe: %w", err)
		}
	}
	if p.settings.Mode > 3 || p.settings.Mode < 0 {
		p.settings.Mode = 0
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	r := Release{}
	var err error
	if p.settings.PreRelease {
		r, err = FetchLatestRelease((*p.settings.baseURL).String())
	} else {
		r, err = FetchLatestStableRelease((*p.settings.baseURL).String())
	}
	if err != nil {
		return fmt.Errorf("failed to check the release: %w", err)
	}
	if len(p.settings.Version) < 1 && len(p.settings.Pipe) < 1 {
		fmt.Println(r.Name)
		return nil
	}

	if r.ReleaseToVersion().EQ(semver.Version{}) {
		if strings.Compare(p.settings.Version, r.Name) != 0 {
			r.Print()
		} else {
			fmt.Printf("Is latest. Your Version %s - Latest: %s\n", p.settings.Version, r.Name)
		}
	} else if r.ReleaseToVersion().GT(p.settings.version) {
		r.Print()
	} else {
		fmt.Printf("Not latest. Your Version %s - Latest: %s\n", p.settings.version.String(), r.Name)
	}
	return nil
}

func gitHubURLs(gh string) (*url.URL, error) {
	baseURL, _ := url.Parse(fmt.Sprintf(BaseURL, gh))
	return baseURL, nil
}
