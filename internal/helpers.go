package internal

import (
	"path/filepath"
)

// RelativeToProject transforms a path to be relative to a base path of the project file
func RelativeToProject(projectFile string, path string) string {
	if projectFile != "" && !filepath.IsAbs(path) {
		projectDir := filepath.Dir(projectFile)
		fullpath := filepath.Join(projectDir, path)
		return fullpath
	}
	return path
}
