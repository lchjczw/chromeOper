package env

import (
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


