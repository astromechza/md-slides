package pdfrenderer

import (
	"os"
	"os/exec"
)

var macOsPaths = []string{
	"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
	"/Applications/Google Chrome Beta.app/Contents/MacOS/Google Chrome",
	"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome",
	"/Applications/Google Chrome Dev.app/Contents/MacOS/Google Chrome",
}

var commandNames = []string{
	"google-chrome-unstable",
	"google-chrome-beta",
	"google-chrome",
	"chromium",
	"chrome",
}

func DetectPathToChrome() string {
	for _, cmd := range commandNames {
		fullPath, err := exec.LookPath(cmd)
		if err == nil {
			return fullPath
		}
	}
	for _, path := range macOsPaths {
		fi, err := os.Stat(path)
		if err == nil && !fi.IsDir() {
			return path
		}
	}
	return ""
}
