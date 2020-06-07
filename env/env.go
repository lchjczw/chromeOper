package env

import (
	"context"
	"errors"
	"github.com/chromedp/chromedp"
)

type EnvInter interface {
	Env() *Env
}
type Env struct {
	Proxy      string
	ChromePath string
}

//将目录添加到path环境变量
func (a *Env) SetChromeExtPath() (chromedp.ExecAllocatorOption, error) {
	//dirTmp, _ := filepath.Split(a.ChromePath)
	//
	//dir, err := filepath.Abs(dirTmp)
	//if err != nil {
	//	return nil, err
	//}

	return chromedp.ExecPath(a.ChromePath), nil
}

//将目录添加到path环境变量
func (a *Env) SetChromeProxy() (chromedp.ExecAllocatorOption, error) {
	if len(a.Proxy) == 0 {
		return nil, errors.New("代理不存在")
	}

	return chromedp.ProxyServer(a.Proxy), nil
}

func (a *Env) Init() []chromedp.ExecAllocatorOption {
	var options []chromedp.ExecAllocatorOption

	defaultOption := chromedp.DefaultExecAllocatorOptions
	options = append(options, defaultOption[:]...)

	//添加路径
	if p, err := a.SetChromeExtPath(); err == nil {
		options = append(options, p)
	}
	//添加代理
	if proxy, err := a.SetChromeProxy(); err == nil {
		options = append(options, proxy)
	}
	return options
}

func NewChrome(env EnvInter) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	e:=env.Env()

	if e == nil {
		panic("初始化错误")
	}

	options := env.Env().Init()

	ctx1, cancel1 := chromedp.NewExecAllocator(ctx, options...)
	// create context
	ctx2, cancel2 := chromedp.NewContext(ctx1)
	cancel := func() {
		cancel1()
		cancel2()
	}

	err := chromedp.Run(ctx2)
	if err != nil {
		defer cancel()
		panic(err)
	}

	return ctx2, cancel

}
