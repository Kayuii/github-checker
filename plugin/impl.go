// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"fmt"
	"net/url"

	"github.com/blang/semver"
)

var BaseURL = "https://api.github.com/repos/%s/%s/releases"

// Settings for the plugin.
type Settings struct {
	GitHubURL  string
	Version    string
	PreRelease bool

	baseURL *url.URL
	version semver.Version
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	var err error

	if p.pipeline.Build.Event != "tag" {
		return fmt.Errorf("github release plugin is only available for tags")
	}

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
	}
	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {

	r := Release{}
	var err error

	if p.settings.PreRelease {
		r, err = FetchLatestRelease(*p.settings.baseURL)
	} else {
		r, err = FetchLatestStableRelease(*p.settings.baseURL)
	}
	if err != nil {
		return fmt.Errorf("failed to check the release: %w", err)
	}

	if len(p.settings.Version) < 1 {
		fmt.Println(r.Name)
		return nil
	}

	if r.Version().GT(p.settings.version) {
		if len(r.Assets) > 0 {
			for _, asset := range r.Assets {
				fmt.Printf("URL: %s\n", asset.URL)
			}
		}
	} else {
		fmt.Printf("Not latest. Your Version %s - Latest: %s\n", p.settings.Version, r.Name)
	}

	return nil
}

func gitHubURLs(gh string) (*url.URL, error) {
	uri, err := url.Parse(gh)
	if err != nil {
		return nil, fmt.Errorf("could not parse GitHub link")
	}

	// Remove the path in the case that DRONE_REPO_LINK was passed in
	uri.Path = ""

	if uri.Hostname() != "github.com" {
		relBaseURL, _ := url.Parse("./api/v3/")
		return uri.ResolveReference(relBaseURL), nil
	}
	baseURL, _ := url.Parse("https://api.github.com/")
	return baseURL, nil
}
