package log

import (
	"fmt"
	"os"
	"runtime"
)

func setDefPath() (string, error) {
	var winDefaultPath = "C:\\Users\\%s\\AppData\\Local\\blacklist\\logs"
	var linuxDefaultPath = ""

	if runtime.GOOS == "windows" {
		var name = os.Getenv("USERNAME")
		return fmt.Sprintf(winDefaultPath, name), nil
	}

	// later
	if runtime.GOOS == "linux" {
		return linuxDefaultPath, nil
	}
	return "", ErrUnknownPlatform
}
