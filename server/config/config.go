package config

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/lchjczw/chromeOper/env"
	"path/filepath"
	"server/util"
)

type Config struct {
	env *env.Env
}

func (a *Config) Env() *env.Env {
	return a.env
}

var (
	global *Config
)

// LoadGlobal 加载全局配置
func LoadGlobal(fpath string) error {
	if global != nil {
		return nil
	}
	c, err := Parse(fpath)
	if err != nil {
		return err
	}
	global = c
	return nil
}

// Global 获取全局配置
func Global() *Config {
	if global == nil {
		return &Config{env: &env.Env{ChromePath:`C:\Users\Administrator\AppData\Local\Google\Chrome\Application\chrome.exe` }}
	}

	return global
}

// Parse 解析配置文件
func Parse(path string) (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func WriteConfig(name string) error {
	if util.IsFile(name) {
		return nil
	}

	dir, _ := filepath.Split(name)
	_ = util.CreateDir(dir)

	buf := new(bytes.Buffer)
	buf.WriteString(`
env:
	proxy: sdfsdf
	path: sdfsdf
`)
	return util.WriteFile(name, buf)
}
