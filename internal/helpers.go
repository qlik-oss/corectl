package internal

import (
	"path/filepath"
)

func RelativeToProject(projectFile string, path string) string {

	if projectFile != "" && !filepath.IsAbs(path) {
		projectDir := filepath.Dir(projectFile)
		fullpath := filepath.Join(projectDir, path)
		return fullpath
	}
	return path
}
