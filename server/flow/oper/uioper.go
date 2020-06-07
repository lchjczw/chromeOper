package oper

import (
	"context"
	"errors"
	"fmt"
	"github.com/lchjczw/chromeOper"
	"net/url"
	"sort"
)

var ReTry int = 3

type UiOper struct {
	ctx      context.Context
	opers    Opers           // 操作列表
	CallBack chromeOper.Call // 全局回调
	Host     string
}

//Len()
func (s UiOper) Len() int {
	return len(s.opers)
}

//Less():步骤由低到高排序
func (s UiOper) Less(i, j int) bool {
	return s.opers[i].Step < s.opers[j].Step
}

//Swap()
func (s UiOper) Swap(i, j int) {
	s.opers[i], s.opers[j] = s.opers[j], s.opers[i]
}

func (a *UiOper) LogErr() {
	fmt.Println("上报错误和截图")
	//	todo 上报错误日志

	a.End()
}
func (a *UiOper) End() {
	//TODO 清理数据
}

func (a *UiOper) run(oper *Oper) (*chromeOper.Oper, error) {
	one_step := oper.ToOper()

	return one_step, chromeOper.RunStep(a.ctx, one_step)

	//one_step := oper.ToOper()
	//for i := 0; i < ReTry; i++ {
	//	err := chromeOper.RunStep(a.ctx, one_step)
	//	if err == nil {
	//		return one_step, nil
	//	}
	//	if i+1 == ReTry {
	//		return nil,err
	//	}
	//}
	//return nil, errors.New("执行失败")
}

func (a *UiOper) Sort() {
	if len(a.opers) > 1 {
		sort.Sort(a)
	}
}

func (a *UiOper) GetOperByOper(oper string) *Oper {
	return a.opers.GetOperByOper(oper)
}
func (a *UiOper) GetOperByName(name string) *Oper {
	return a.opers.GetOperByName(name)
}

func (a *UiOper) RepleaceHost() {
	if a.Host == "" {
		return
	}
	for _, v := range a.opers {
		if v.Oper != "OpenUrl" {
			continue
		}
		u, err := url.Parse(v.Args)
		if err != nil {
			continue
		}
		u.Host = a.Host
		v.Args = u.String()
	}
}

//执行oper
func (a *UiOper) Run() error {
	if len(a.opers) == 0 || a.opers == nil {
		return errors.New("请添加执行操作")
	}
	defer a.ClearOper()

	a.RepleaceHost()

	for _, oper := range a.opers {

		//回调处理
		for i := 0; i < ReTry; i++ {
			//执行操作步骤
			res, err := a.run(oper)
			if err != nil {
				fmt.Println(err.Error())
				a.LogErr()
				// 持续运行到最后一次，还是错误则直接返回，否则重复运行，直到正确
				if i+1 == ReTry {
					return err
				}
				continue
				//return errors.New("执行中断")
			}
			if a.CallBack == nil {
				break
			}
			err = a.CallBack(a.ctx, res)
			//回调处理错误，重复三次
			if err == nil {
				break
			}
			fmt.Println(err.Error())
		}
	}

	//成功
	return nil

}

//创建oper
func NewUiOper(ctx context.Context, opers Opers, callBack chromeOper.Call) *UiOper {
	oper := new(UiOper)
	oper.ctx = ctx
	if opers == nil {
		var ops Opers
		opers = ops
	}
	oper.opers = opers
	oper.CallBack = callBack
	return oper
}

func (a *UiOper) ReFresh() *UiOper {
	a.ClearOper()
	a.Add("刷新页面", "Reload", ``, "", 1, 500)
	a.Run()
	return a
}

func (a *UiOper) ClearOper() *UiOper {
	var opers Opers
	a.opers = opers
	return a
}

//添加步骤
func (a *UiOper) Add(name, oper, sel string, args string, step, timeOut int) *UiOper {
	index := len(a.opers)
	if step == 0 {
		if index == 0 {
			step = 5
		} else {
			step = a.opers[index-1].Step + 5
		}
	}
	o := NewOper(name, oper, sel, step, args, timeOut)
	if o != nil {
		a.opers = append(a.opers, o)
	}
	return a
}
func (a *UiOper) AddOpers(ops Opers) {
	a.opers = append(a.opers, ops...)
}
