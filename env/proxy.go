package env

import (
	"github.com/chromedp/chromedp"
)

//将目录添加到path环境变量
func SetChromeProxy(env EnvInter) (chromedp.ExecAllocatorOption, error) {
	if len(env.ChromeProxy()) == 0 {
		return nil, nil
	}

	return chromedp.ProxyServer(env.ChromeProxy()), nil
}
