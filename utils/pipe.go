package utils

import (
	"bufio"
	"os"
	"strings"
)

func GetPipe() (string, bool) {
	fileInfo, _ := os.Stdin.Stat()
	if (fileInfo.Mode() & os.ModeNamedPipe) != os.ModeNamedPipe {
		return "", false
	}
	s := bufio.NewScanner(os.Stdin)
	resList := make([]string, 0)
	for s.Scan() {
		resList = append(resList, s.Text())
	}
	result := strings.Join(resList, "\n")
	return result, true
}
