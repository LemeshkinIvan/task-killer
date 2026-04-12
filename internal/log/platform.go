package log

import (
	"runtime"
)

func setDefPath() (string, error) {
	if runtime.GOOS == "windows" {
		return winDefaultPath, nil
	}

	// if runtime.GOOS == "linux" {

	// }
	return "", ErrUnknownPlatform
}

var winDefaultPath = "C:\\Users\\%s\\AppData\\Local\\blacklist\\logs"
