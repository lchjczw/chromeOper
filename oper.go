package ui

import (
	ck "chromeOper/cookies"
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//截图数量
var num, min_num, max_num int = 0, 0, 20
//是否开启debug
var debug bool = false
var img_dir string = "./img/"

func getDevice(dev string) chromedp.Device {
	var devi chromedp.Device = nil
	switch dev {
	case "IPad":
		devi = device.IPad
		break
	case "IPhone8":
		devi = device.IPhone8
		break
	default:
		devi = device.IPhoneX
		break
	}

	return devi

}


//点击
func ClickTime(ctx context.Context, sel string, t time.Duration) error {
	var ui []byte
	err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&ui),
		chromedp.Click(sel),
		chromedp.Sleep(t),
	)
	WriteImg(ui, "ClickTime")

	return err
}
func Click(ctx context.Context, sel string) error {
	var ui []byte
	err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&ui),
		chromedp.Click(sel),
	)
	WriteImg(ui, "Click")

	return err
}
func Submit(ctx context.Context, sel string) error {
	var ui []byte
	err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&ui),
		chromedp.Submit(sel, chromedp.NodeVisible),
	)
	WriteImg(ui, "Submit")

	return err
}
func ClickByQuery(ctx context.Context, sel string) error {
	var ui []byte
	err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&ui),
		chromedp.Click(sel, chromedp.ByQuery),
	)
	WriteImg(ui, "ClickByQuery")

	return err
}

func ClickByQueryTime(ctx context.Context, sel string, t time.Duration) error {
	var ui []byte
	err := chromedp.Run(ctx,
		//chromedp.WaitVisible(sel, chromedp.ByQuery),
		chromedp.Sleep(t),
		chromedp.CaptureScreenshot(&ui),
		chromedp.Click(sel, chromedp.ByQuery),
	)
	WriteImg(ui, "ClickByQueryTime")

	return err
}
func SetDevice(ctx context.Context, device string) error {
	dev := getDevice(device)
	return chromedp.Run(ctx, chromedp.Emulate(dev))
}
func OpenUrl(ctx context.Context, url string) error {
	var ui []byte
	err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		chromedp.CaptureScreenshot(&ui),
	)
	WriteImg(ui, "OpenUrl")

	return err
}
func Sleep(ctx context.Context, t time.Duration) error {
	var ui []byte

	err := chromedp.Run(
		ctx,
		chromedp.Sleep(t),
		chromedp.CaptureScreenshot(&ui),
	)

	WriteImg(ui, "Sleep")

	return err
}
func Reload(ctx context.Context, t time.Duration) error {
	var ui []byte

	err := chromedp.Run(
		ctx,
		chromedp.Reload(),
		chromedp.Sleep(t),
		chromedp.CaptureScreenshot(&ui),

	)

	WriteImg(ui, "Reload")

	return err
}
func SendKeys(ctx context.Context, sel, val string) error {
	var ui []byte

	err := chromedp.Run(
		ctx,
		chromedp.SendKeys(sel, val, chromedp.NodeVisible),
		chromedp.CaptureScreenshot(&ui),

	)

	WriteImg(ui, "SendKeys")

	return err
}
func SetValue(ctx context.Context, sel, val string) error {
	var ui []byte

	err := chromedp.Run(
		ctx,
		chromedp.SetValue(sel, val, chromedp.NodeVisible),
		chromedp.CaptureScreenshot(&ui),
	)

	WriteImg(ui, "SetValue")

	return err
}
func Capture(ctx context.Context, sel, logs string) error {
	var ui []byte

	var err error
	if sel != "" {
		err = chromedp.Run(
			ctx,
			chromedp.Screenshot(sel, &ui),
		)
	} else {
		err = chromedp.Run(
			ctx,
			chromedp.CaptureScreenshot(&ui),
		)
	}

	WriteImg(ui, "Capture")

	log.Println("截图：", fmt.Sprintf("%03d", num)+".Capture.png", logs)

	return err
}
func ClickByQueryWaitNoVisible(ctx context.Context, sel string) error {
	var ui []byte
	err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&ui),
		chromedp.Click(sel, chromedp.ByQuery),
		chromedp.WaitNotVisible(sel, chromedp.ByQuery),
	)

	WriteImg(ui, "ClickByQueryWaitNoVisible")

	return err
}
func ClickWaitNoVisible(ctx context.Context, sel string) error {
	var ui []byte
	err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&ui),
		chromedp.Click(sel),
		chromedp.WaitNotVisible(sel),
	)
	WriteImg(ui, "ClickWaitNoVisible")

	return err
}
func GetText(ctx context.Context, sel string, v *string) error {
	var ui []byte
	err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&ui),
		chromedp.Text(sel, v),
	)

	WriteImg(ui, "GetText")
	return err
}

func GetOuterHTML(ctx context.Context, sel string, v *string) error {
	var ui []byte
	err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&ui),
		chromedp.OuterHTML(sel, v),
	)

	WriteImg(ui, "GetOuterHTML")
	return err
}
func GetValue(ctx context.Context, sel string, v *string) error {
	var ui []byte
	err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&ui),
		chromedp.Value(sel, v),
	)

	WriteImg(ui, "GetValue")
	return err
}
func GetAttributeValue(ctx context.Context, sel, name string, v *string, ok *bool) error {
	var ui []byte
	err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&ui),
		chromedp.AttributeValue(sel, name, v, ok),
	)
	WriteImg(ui, "GetAttributeValue")

	return err
}

func ClickLoopTime(ctx context.Context, sel, count string, t time.Duration) error {
	var ui []byte
	var err error = nil
	n, _ := strconv.Atoi(count)

	for i := 0; i < n; i++ {
		err = chromedp.Run(ctx,
			chromedp.CaptureScreenshot(&ui),
			chromedp.Click(sel),
			chromedp.Sleep(t),
		)
		if err != nil {
			break
		}
	}

	WriteImg(ui, "ClickLoopTime")
	return err
}

// BetUiOper bet_ui_oper对象
type Oper struct {
	Name     string //操作标识，同一个gettext，可以不同的标识，以执行不同的逻辑
	Value    string //参数 ,value body
	Oper     string //执行操作 ,oper url
	Sel      string //界面元素定位,sel header
	TimeOut  string //time_out
	OtherArg interface{}
	Result
	Hook
}

type PreCall func(ctx context.Context, oper *Oper)
type AfterCall func(ctx context.Context, oper *Oper)

type Hook struct {
	Pre   PreCall
	After AfterCall
}

type Result struct {
	Res      string
	Ok       bool
	OtherRes *interface{}
	//SuccessOper *Oper //这个oper必须是api，提交数据用的
	//FailOper *Oper //这个oper必须是api，提交数据用的
}

//url,header,body
func (a *Oper) GetOperType() string {
	if strings.HasPrefix(a.Oper, "http") {
		return "url"
	}

	return a.Oper
}

//url,header,body

//请求单条数据 做简单的通知功能
//成功返回body 失败返回空字符串
func (a *Oper) Request() error {
	//判断是否需要请求数据
	if a.GetOperType() != "url" {
		return errors.New("非api操作")
	}
	request := gorequest.New()
	AddHeader(request, a.Sel)
	url := a.Oper
	//存在body则是post，不存在则是get
	if a.Value != "" {
		request = request.Post(url).Send(a.Value)
	} else {
		request = request.Get(url)
	}

	resp, body, errs := request.End()

	if resp.StatusCode != http.StatusOK || len(errs) > 0 {
		return errors.New("请求失败")
	}

	//	解析body数据
	fmt.Println(body)

	return nil
}

func UiRefresh(ctx context.Context) error {

	oper := new(Oper)
	oper.TimeOut = "30"
	oper.Oper = "Reload"

	return RunStep(ctx, oper)
}

func UiSleep(ctx context.Context, t int) error {

	oper := new(Oper)
	oper.TimeOut = fmt.Sprintf("%d", t)
	oper.Oper = "Sleep"

	return RunStep(ctx, oper)

}

func RunStep(ctx context.Context, oper *Oper) error {

	//执行pre hook
	if oper.Pre != nil {
		oper.Pre(ctx, oper)
	}

	var err error = nil
	switch oper.GetOperType() {
	case "SetValue":
		err = SetValue(ctx, oper.Sel, oper.Value)
		break
	case "SendKeys":
		err = SendKeys(ctx, oper.Sel, oper.Value)
		break
	case "Reload":
		err = Reload(ctx, transTime(oper.TimeOut))
		break
	case "Sleep":
		err = Sleep(ctx, transTime(oper.TimeOut))
		break
	case "OpenUrl":
		err = OpenUrl(ctx, oper.Value)
		break
	case "SetDevice":
		err = SetDevice(ctx, oper.Value)
		break
	case "ClickByQueryTime":
		err = ClickByQueryTime(ctx, oper.Sel, transTime(oper.TimeOut))
		break
	case "ClickByQuery":
		err = ClickByQuery(ctx, oper.Sel)
		break
	case "Click":
		err = Click(ctx, oper.Sel)
		break
	case "Capture":
		err = Capture(ctx, oper.Sel, oper.Value)
		break
	case "ClickTime":
		err = ClickTime(ctx, oper.Sel, transTime(oper.TimeOut))
		break
	case "ClickByQueryWaitNoVisible":
		err = ClickByQueryWaitNoVisible(ctx, oper.Sel)
		break
	case "ClickWaitNoVisible":
		err = ClickWaitNoVisible(ctx, oper.Sel)
		break
	case "ClickLoopTime":
		err = ClickLoopTime(ctx, oper.Sel, oper.Value, transTime(oper.TimeOut))
		break
	case "url":
		err = oper.Request()
		break
	case "GetText":
		err = GetText(ctx, oper.Sel, &oper.Res)
		break
	case "GetOuterHTML":
		err = GetOuterHTML(ctx, oper.Sel, &oper.Res)
		break
	case "GetValue":
		err = GetValue(ctx, oper.Sel, &oper.Res)
		break
	case "GetAttributeValue":
		err = GetAttributeValue(ctx, oper.Sel, oper.Value, &oper.Res, &oper.Ok)
		break
	case "Submit":
		err = Submit(ctx, oper.Sel)
		break
	case "GetCookies":
		err = GetCookies(ctx, oper.OtherRes)
		break
	case "SetCookies":
		err = SetCookies(ctx, oper.OtherArg)
		break

	default:
		err = errors.New("不支持的操作")
		break
	}

	if err == nil {
		oper.Ok = true
	} else {
		oper.Ok = false
	}

	log.Println(num, oper.Name, oper.Oper, oper.Sel, oper.Value, oper.TimeOut, oper.Ok, strings.Fields(oper.Res))

	//执行after hook
	if oper.After != nil {
		oper.After(ctx, oper)
	}

	return err
}

func GetCookies(ctx context.Context, Cookies *interface{}) error {

	var ui []byte
	err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&ui),
		chromedp.ActionFunc(ck.GetChromedpCookies),
	)

	*Cookies = *ck.GetCookies()

	WriteImg(ui, "GetCookies")

	return err
}

func SetCookies(ctx context.Context, Cookies interface{}) error {

	ck.SetCookies(Cookies.([]*network.Cookie))

	var ui []byte
	err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&ui),
		chromedp.ActionFunc(ck.SetChromedpCookies),
	)

	WriteImg(ui, "SetCookies")

	return err
}


func NewChromeTab(ctx context.Context) (context.Context, context.CancelFunc) {

	err := CheckEnv()
	if err != nil {
		panic(err)
	}

	ctx, cancel := chromedp.NewContext(
		ctx,
		//chromedp.WithLogf(log.Printf),
		//chromedp.WithDebugf(log.Printf),
	)

	err = chromedp.Run(ctx)
	if err != nil {
		panic(err)
	}

	return ctx, cancel
}

func AddTimeOut(ctx context.Context, t int) (context.Context, context.CancelFunc) {
	//超时设置
	return context.WithTimeout(ctx, time.Duration(t)*time.Second)
}

func NewChrome(ctx context.Context, is_show bool) (context.Context, context.CancelFunc) {

	err := CheckEnv()
	if err != nil {
		panic(err)
	}

	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		panic(err)
	}

	var opts []chromedp.ExecAllocatorOption
	if is_show {
		opts = append(chromedp.DefaultExecAllocatorOptions[0:2], chromedp.DefaultExecAllocatorOptions[4:]...)

		opts = append(opts,
			//chromedp.DisableGPU,
			chromedp.UserDataDir(dir),
		)
	} else {
		opts = append(opts, chromedp.DefaultExecAllocatorOptions[:]...)
	}

	allocCtx, cancel1 := chromedp.NewExecAllocator(ctx, opts...)
	//defer cancel()

	// also set up a custom logger
	taskCtx, cancel2 := chromedp.NewContext(allocCtx)

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}

	return taskCtx, func() {
		defer cancel1()
		defer cancel2()
		defer os.RemoveAll(dir)
	}

}
