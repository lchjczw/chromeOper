package flow

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/lchjczw/chromeOper"
	"log"
	"server/flow/oper"
	"strings"
)


func RunChromeDp(c *gin.Context, ctx context.Context) {
	oper := new(chromeOper.Oper)

	err := c.BindJSON(oper)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
	}

	err = chromeOper.RunStep(ctx, oper)
	if err != nil {
		return
	}
}

func Login(ctx context.Context, userName, password string) (res string) {
	//oper := new(chromeOper.Oper)

	uiOper := oper.NewUiOper(ctx, nil, nil)
	// uiOper.Add("设置设备", "SetDevice", "", "IPhoneX", 0, 0)
	uiOper.Add("等待", "Sleep", "", "3000", 0, 0)
	uiOper.Add("打开网页", "OpenUrl", "", "https://www.battlenet.com.cn/login/zh/", 0, 0)
	uiOper.Add("等待", "Sleep", "", "3000", 0, 0)



	//uiOper.Add("获取验证码", "Capture", "#sec-string", "verify", 0, 0)
	//if err := uiOper.Run(); err != nil {
	//	log.Println(err.Error())
	//}
	//
	//uiOper.ClearOper()
	//
	//util.ShowVerify()
	//verify := util.InputVerify()
	//if verify != "" {
	//	uiOper.Add("设置验证码", "SetValue", "#captchaInput", verify, 0, 0)
	//}

	uiOper.Add("设置昵称", "SetValue", "#accountName", userName, 0, 0)
	uiOper.Add("设置密码", "SetValue", "#password", password, 0, 0)
	uiOper.Add("登录", "Click", "#submit", "", 0, 0)
	uiOper.Add("反馈信息", "GetText", "#display-errors", "", 0, 0) // 找不到该暴雪游戏通行证。

	uiOper.CallBack = func(ctx context.Context, oper *chromeOper.Oper) error {
		if oper.Name == "反馈信息" {
			log.Println("结果:", oper.Res)
			r:=strings.Split(oper.Res,"\n")
			if len(r) > 0 {
				res = r[0]
			}
		}
		return nil
	}

	err := uiOper.Run()
	if err != nil {
		log.Println(err.Error())
	}
	if res != "" {
		if strings.Contains(res, "找不到该暴雪游戏通行证") {
			res = "账号不存在"
		}
		if strings.Contains(res, "密码") {
			res = "密码不对"
		}
		if strings.Contains(res, "包含错误字符") {
			res = "验证码不对"
		}

	}

	//chromedp.Run(ctx,chromedp.Query(`sdfssdfsdf`,chromedp.BySearch))
	return
}
