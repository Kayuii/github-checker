// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"io/ioutil"
	"os"
	"regexp"

	"github.com/blang/semver"
)

var sem = regexp.MustCompile(`(?:^|\D)(\d+\.\d+\.\d+\S*)(?:$|\s)`)

func readStringOrFile(input string) (string, error) {
	if len(input) > 255 {
		return input, nil
	}
	// Check if input is a file path
	if _, err := os.Stat(input); err != nil && os.IsNotExist(err) {
		// No file found => use input as result
		return input, nil
	} else if err != nil {
		return "", err
	}
	result, err := ioutil.ReadFile(input)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func Version(input string) semver.Version {
	match := sem.FindStringSubmatch(input)
	if len(match) < 2 {
		return semver.Version{}
	}
	if sv, err := semver.ParseTolerant(match[1]); err == nil {
		return sv
	}
	return semver.Version{}
}
