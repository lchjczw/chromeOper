package chromeOper

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	ck "github.com/lchjczw/chromeOper/cookies"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

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
	err := chromedp.Run(ctx,
		chromedp.Click(sel),
		chromedp.Sleep(t),
	)

	WriteImg(ctx, "ClickTime")

	return err
}
func Click(ctx context.Context, sel string) error {
	err := chromedp.Run(ctx,
		chromedp.Click(sel),
	)
	WriteImg(ctx, "Click")

	return err
}
func Submit(ctx context.Context, sel string) error {
	err := chromedp.Run(ctx,
		chromedp.Submit(sel, chromedp.NodeVisible),
	)
	WriteImg(ctx, "Submit")

	return err
}
func ClickByQuery(ctx context.Context, sel string) error {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible(sel,chromedp.ByQuery),
		chromedp.Click(sel, chromedp.ByQuery),
	)
	WriteImg(ctx, "ClickByQuery")

	return err
}

func ClickByQueryTime(ctx context.Context, sel string, t time.Duration) error {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible(sel,chromedp.ByQuery),
		chromedp.Click(sel, chromedp.ByQuery),
		chromedp.Sleep(t),
	)
	WriteImg(ctx, "ClickByQueryTime")

	return err
}
func SetDevice(ctx context.Context, device string) error {
	dev := getDevice(device)
	return chromedp.Run(ctx, chromedp.Emulate(dev))
}
func OpenUrl(ctx context.Context, url string) error {
	err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),
	)
	WriteImg(ctx, "OpenUrl")

	return err
}
func Sleep(ctx context.Context, t time.Duration) error {
	//time.Sleep(t)
	return chromedp.Run(
		ctx,
		chromedp.Sleep(t),
	)
}

func Reload(ctx context.Context, t time.Duration) error {
	err := chromedp.Run(
		ctx,
		chromedp.Reload(),
	)

	Sleep(ctx, t)
	WriteImg(ctx, "Reload")

	return err
}
func SendKeys(ctx context.Context, sel, val string) error {
	err := chromedp.Run(
		ctx,
		chromedp.SendKeys(sel, val, chromedp.NodeVisible),
	)

	WriteImg(ctx, "SendKeys")

	return err
}
func SetValue(ctx context.Context, sel, val string) error {
	err := chromedp.Run(
		ctx,
		chromedp.SetValue(sel, val, chromedp.NodeVisible),
	)

	WriteImg(ctx, "SetValue")

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

	SaveImage(ctx, ui, "Capture")

	log.Println("截图：", fmt.Sprintf("%03d", num)+".Capture.png", logs)

	return err
}
func ClickByQueryWaitNoVisible(ctx context.Context, sel string) error {
	err := chromedp.Run(ctx,
		chromedp.Click(sel, chromedp.ByQuery),
		chromedp.WaitNotVisible(sel, chromedp.ByQuery),
	)

	WriteImg(ctx, "ClickByQueryWaitNoVisible")

	return err
}
func ClickWaitNoVisible(ctx context.Context, sel string) error {
	err := chromedp.Run(ctx,
		chromedp.Click(sel),
		chromedp.WaitNotVisible(sel),
	)
	WriteImg(ctx, "ClickWaitNoVisible")

	return err
}
func GetText(ctx context.Context, sel string, v *string) error {
	err := chromedp.Run(ctx,
		chromedp.Text(sel, v),
	)

	WriteImg(ctx, "GetText")
	return err
}

func GetOuterHTML(ctx context.Context, sel string, v *string) error {
	err := chromedp.Run(ctx,
		chromedp.OuterHTML(sel, v),
	)

	WriteImg(ctx, "GetOuterHTML")
	return err
}
func GetValue(ctx context.Context, sel string, v *string) error {

	err := chromedp.Run(ctx,
		chromedp.Value(sel, v, chromedp.NodeVisible),
	)

	WriteImg(ctx, "GetValue")
	return err
}
func GetAttributeValue(ctx context.Context, sel, name string, v *string, ok *bool) error {
	err := chromedp.Run(ctx,
		chromedp.AttributeValue(sel, name, v, ok, chromedp.NodeVisible),
	)
	WriteImg(ctx, "GetAttributeValue")

	return err
}

func ClickLoopTime(ctx context.Context, sel, count string, t time.Duration) error {
	var err error = nil
	n, _ := strconv.Atoi(count)

	for i := 0; i < n; i++ {
		err = chromedp.Run(ctx,
			chromedp.Click(sel),
			chromedp.Sleep(t),
		)
		WriteImg(ctx, "ClickLoopTime")
		if err != nil {
			break
		}
	}


	return err
}

// BetUiOper bet_ui_oper对象
type Oper struct {
	Name     string        //操作标识，同一个gettext，可以不同的标识，以执行不同的逻辑
	Value    string        //参数 ,value body
	Oper     string        //执行操作 ,oper url
	Sel      string        //界面元素定位,sel header
	TimeOut  time.Duration //time_out
	OtherArg interface{}
	Result
	Hook
}

type Call func(ctx context.Context, oper *Oper) error

type Hook struct {
	Pre   Call
	After Call
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
	oper.TimeOut = 30
	oper.Oper = "Reload"

	return RunStep(ctx, oper)
}

func UiSleep(ctx context.Context, t time.Duration) error {

	oper := new(Oper)
	oper.TimeOut = t
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
		err = Reload(ctx, oper.TimeOut)
		break
	case "Sleep":
		err = Sleep(ctx, oper.TimeOut)
		break
	case "OpenUrl":
		err = OpenUrl(ctx, oper.Value)
		break
	case "SetDevice":
		err = SetDevice(ctx, oper.Value)
		break
	case "ClickByQueryTime":
		err = ClickByQueryTime(ctx, oper.Sel, oper.TimeOut)
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
		err = ClickTime(ctx, oper.Sel, oper.TimeOut)
		break
	case "ClickByQueryWaitNoVisible":
		err = ClickByQueryWaitNoVisible(ctx, oper.Sel)
		break
	case "ClickWaitNoVisible":
		err = ClickWaitNoVisible(ctx, oper.Sel)
		break
	case "ClickLoopTime":
		err = ClickLoopTime(ctx, oper.Sel, oper.Value, oper.TimeOut)
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

	err := chromedp.Run(ctx,
		chromedp.ActionFunc(ck.GetChromedpCookies),
	)

	*Cookies = *ck.GetCookies()

	WriteImg(ctx, "GetCookies")

	return err
}

func SetCookies(ctx context.Context, Cookies interface{}) error {

	ck.SetCookies(Cookies.([]*network.Cookie))

	err := chromedp.Run(ctx,
		chromedp.ActionFunc(ck.SetChromedpCookies),
	)

	WriteImg(ctx, "SetCookies")

	return err
}
func NewChromeDp() (context.Context, context.CancelFunc) {

	err := CheckEnv()
	if err != nil {
		panic(err)
	}

	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)

	err = chromedp.Run(ctx)
	if err != nil {
		panic(err)
	}

	return ctx, cancel
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

func AddTimeOut(ctx context.Context, t time.Duration) (context.Context, context.CancelFunc) {
	//超时设置
	return context.WithTimeout(ctx, t)
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
