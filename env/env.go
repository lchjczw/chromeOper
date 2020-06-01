package env

import (
	"context"
	"github.com/chromedp/chromedp"
)

type EnvInter interface {
	ChromeProxy() string
	ChromeExecPath() string
}

func NewChrome(env EnvInter) (context.Context, context.CancelFunc) {
	ctx := context.Background()

	var options []chromedp.ExecAllocatorOption

	defaultOption := chromedp.DefaultExecAllocatorOptions
	options = append(options, defaultOption[:]...)

	//添加路径
	if p, err := SetChromeExtPath(env); err == nil {
		options = append(options, p)
	}
	//添加代理
	if proxy, err := SetChromeProxy(env); err == nil {
		options = append(options, proxy)
	}

	ctx1, cancel1 := chromedp.NewExecAllocator(ctx, options...)
	// create context
	ctx2, cancel2 := chromedp.NewContext(ctx1)
	cancel:= func(){
		cancel1()
		cancel2()
	}

	err := chromedp.Run(ctx)
	if err != nil {
		defer cancel()
		panic(err)
	}

	return ctx2, cancel

}
