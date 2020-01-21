package chromeOper

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

//截图数量
var num, min_num, max_num int = 0, 0, 20
//是否开启debug
var debug bool = true
var is_sync bool = false
var img_dir string = "./img/"

func WriteImg(ctx context.Context, logs string) {

	//声明式调试时，才截图
	if !debug {
		return
	}

	//不存在则创建目录
	if !PathExists(img_dir) {
		//	创建目录
		if err := os.Mkdir(img_dir, os.ModePerm); err != nil {
			return
		}
	}

	var ui []byte
	if is_sync {
		go SaveImage(ctx, ui, logs)
	}else {
		SaveImage(ctx, ui, logs)
	}

	return
}

//保存截图
func SaveImage(ctx context.Context, ui []byte, logs string) {

	if len(ui) == 0 {
		err := chromedp.Run(ctx,
			chromedp.CaptureScreenshot(&ui),
		)
		if err != nil {
			return
		}
	}

	if num == max_num {
		num = min_num
	}
	//num为全局变量
	num++

	prex := fmt.Sprintf("%03d", num)
	file_name := fmt.Sprintf("%s/%s.%s.png", img_dir, prex, logs)
	//绝对路径
	file_name, err := filepath.Abs(file_name)
	if err != nil {
		log.Println(err)
		return
	}
	//删除对应数字编号的图片
	err = DeleteFileByPreFixName(img_dir, prex);
	if err != nil {
		log.Println(err)
	}
	//保存截图
	err = ioutil.WriteFile(file_name, ui, 0777)
	if err != nil {
		log.Println(err)
	}

	return
}

//删除指定前缀的文件
func DeleteFileByPreFixName(fileDir, preFix string) error {
	files, _ := ioutil.ReadDir(fileDir)
	var err error

	for _, onefile := range files {
		if onefile.IsDir() {
			continue
		}

		if !strings.HasPrefix(onefile.Name(), preFix) {
			continue
		}
		//	执行删除
		if err = os.Remove(fileDir + "/" + onefile.Name()); err != nil {
			return err
		}
	}

	return nil
}

//文件或者目录是否存在
func PathExists(path string) (bool) {
	_, err := os.Stat(path);
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//添加头部数据
func AddHeader(req *gorequest.SuperAgent, head string) *gorequest.SuperAgent {
	if head == "" {
		return req
	}

	header := strings.Split(head, "|")
	for _, value := range header {
		val := strings.Split(value, "=")
		if val[0] != "" && val[1] != "" {
			req = req.Set(val[0], val[1])
		}
	}

	return req
}

//将目录添加到path环境变量
func AddToPath(dir string) error {

	dir, err := filepath.Abs(dir)
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

	err = os.Setenv("PATH", dir+seq+path)
	if err != nil {
		return err
	}

	return nil
}

//复制文件
func CopyFile(src, dst string, is_over_write bool) error {

	if !PathExists(src) {
		return errors.New("源文件不存在")
	}

	//已经存在，不覆盖
	if PathExists(dst) && is_over_write == false {
		return nil
	}

	src_f, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, src_f, 0777)
	if err != nil {
		return err
	}

	return nil
}

func SplitNum(str string, seq string) []int {
	a := strings.Split(str, seq)
	var r []int
	for _, value := range a {
		i, err := strconv.Atoi(value)
		if err != nil {
			continue
		}
		r = append(r, i)
	}

	return r
}
