package internal

import (
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
)

func FixPath(path string) string {
	path = strings.Replace(path, ":\\", "/", -1)
	path = strings.Replace(path, "\\", "/", -1)
	if path[0] != '/' {
		path = "/" + path
	}
	return path
}

func MaybeToWindowsPath(path string) string {
	log.Debug(runtime.GOOS)
	return path
}
