package internal

import "strings"

func FixPath(path string) string {
	path = strings.Replace(path, ":\\", "/", -1)
	path = strings.Replace(path, "\\", "/", -1)
	if path[0] != '/' {
		path = "/" + path
	}
	return path
}
