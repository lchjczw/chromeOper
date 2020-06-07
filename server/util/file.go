package util

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	path, _ = filepath.Abs(path)

	s, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}

	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(file string) bool {
	file, _ = filepath.Abs(file)
	s, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		return false
	}

	return !s.IsDir()
}

// 判断所给路径文件/文件夹是否存在
func IsExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

const (
	ENDLINE = "\n"
	TAP     = "\t"
)

func CreateDir(fileDir string) error {
	if IsFile(fileDir) {
		return errors.New("存在相同的文件名")
	}

	//创建目录
	err := os.MkdirAll(fileDir, 0751)
	if err != nil {
		return errors.New("目录创建失败")
	}

	return nil
}

func ReadDir(fileDir string) ([]os.FileInfo, error) {

	if !IsExists(fileDir) {
		return nil, errors.New("文件夹不存在")
	}

	//创建目录
	return ioutil.ReadDir(fileDir)
}

// 读取文件内容
func ReadFile(name string) (*bytes.Buffer, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	return buf, err
}

// 写入文件
func WriteFile(name string, buf *bytes.Buffer) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, buf)
	if err != nil {
		return err
	}

	return nil
}

