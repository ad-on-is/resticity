package internal

import (
	"runtime"
	"strings"
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
	if runtime.GOOS != "windows" {
		return path
	}
	p := strings.Split(path, "/")
	p = p[1:] // skip first empty string
	p[0] = p[0] + ":"
	path = strings.Join(p, "\\")
	return path
}
