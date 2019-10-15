package ui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func SetRunMod(is_debug bool) {
	debug = is_debug
}

//如果是Windows系统，则添加一下路径
func Init() {
	//if runtime.GOOS == "windows" {
	AddToPath(fmt.Sprintf(`C:\Users\%s\AppData\Local\Google\Chrome\Application\`, os.Getenv("USERNAME")))
	//}

}

func CheckEnv() error {
	if FindChromePath() == "" {
		return errors.New("没找到浏览器路径")
	}
	return nil
}

func FindChromePath() string {
	for _, path := range [...]string{
		// Unix-like
		"headless_shell",
		"headless-shell",
		"chromium",
		"chromium-browser",
		"google-chrome",
		"google-chrome-stable",
		"google-chrome-beta",
		"google-chrome-unstable",
		"/usr/bin/google-chrome",

		// Windows
		"chrome",
		"chrome.exe", // in case PATHEXT is misconfigured
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		fmt.Sprintf(`C:\Users\%s\AppData\Local\Google\Chrome\Application\`, os.Getenv("USERNAME")),
		// Mac
		`/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`,
	} {
		found, err := exec.LookPath(path)
		if err == nil {
			return found
		}
	}
	// Fall back to something simple and sensible, to give a useful error
	// message.
	return ""
}
