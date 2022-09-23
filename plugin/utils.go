// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/blang/semver/v4"
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
	match := sem.FindStringSubmatch(RemoveBrackets(input))
	if len(match) < 2 {
		return semver.Version{}
	}
	if sv, err := semver.ParseTolerant(match[1]); err == nil {
		return sv
	}
	return semver.Version{}
}

func RemoveBrackets(input string) string {
	inputName := ReplaceToBytes(input, " ", "", -1)
	inputName = ReplaceToBytes(string(inputName), "-", ".", -1)
	inputName = ReplaceToBytes(string(inputName), "(", ".", -1)
	inputName = ReplaceToBytes(string(inputName), ")", "", -1)
	return string(inputName)
}

func ReplaceToBytes(s, old, new string, n int) []byte {
	if old == new || n == 0 {
		return []byte(s) // avoid allocation
	}

	// Compute number of replacements.
	if m := strings.Count(s, old); m == 0 {
		return []byte(s) // avoid allocation
	} else if n < 0 || m < n {
		n = m
	}

	// Apply replacements to buffer.
	t := make([]byte, len(s)+n*(len(new)-len(old)))
	w := 0
	start := 0
	for i := 0; i < n; i++ {
		j := start
		if len(old) == 0 {
			if i > 0 {
				_, wid := utf8.DecodeRuneInString(s[start:])
				j += wid
			}
		} else {
			j += strings.Index(s[start:], old)
		}
		w += copy(t[w:], s[start:j])
		w += copy(t[w:], new)
		start = j + len(old)
	}
	w += copy(t[w:], s[start:])
	return t[0:w]
}
