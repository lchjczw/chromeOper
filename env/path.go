package env

import (
	"github.com/chromedp/chromedp"
	"os"
	"path/filepath"
	"runtime"
)

//将目录添加到path环境变量
func AddDirToPath(dir string) error {
	dirTmp, _ := filepath.Split(dir)

	dir, err := filepath.Abs(dirTmp)
	if err != nil {
		return err
	}

	path := os.Getenv("PATH")

	seq := ";"

	if runtime.GOOS == "windows" {
		seq = ";"
	} else {
		seq = ":"
	}

	return os.Setenv("PATH", dir+seq+path)

}

//将目录添加到path环境变量
func SetChromeExtPath(env EnvInter) (chromedp.ExecAllocatorOption, error) {
	dirTmp, _ := filepath.Split(env.ChromeExecPath())

	dir, err := filepath.Abs(dirTmp)
	if err != nil {
		return nil, err
	}

	return chromedp.ExecPath(dir), nil
}
