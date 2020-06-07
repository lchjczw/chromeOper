package oper

import (
	"errors"
	"github.com/lchjczw/chromeOper"
)

// Oper对象
type Oper struct {
	Name        string `json:"name"  binding:"required"`
	Oper        string `json:"oper"  binding:"required"`
	Sel         string `json:"sel"  binding:"required"`
	Args        string `json:"args"`
	Step        int    `json:"step"  binding:"required"`
	TimeOut     int    `json:"time_out"`
	Prev, After chromeOper.Call
}

type Opers []*Oper

func (a Opers) GetOperByName(name string) *Oper {
	for _, oper := range a {
		if oper.Name == name {
			return oper
		}
	}
	return nil
}

func (a Opers) GetOperByOper(oper string) *Oper {
	for _, o := range a {
		if o.Oper == oper {
			return o
		}
	}
	return nil
}

func (a *Oper) ToOper() *chromeOper.Oper {
	one_step := new(chromeOper.Oper)
	one_step.Oper = a.Oper
	one_step.Value = a.Args
	one_step.Sel = a.Sel
	one_step.TimeOut = chromeOper.GetTime(a.TimeOut)
	one_step.Name = a.Name
	one_step.Pre = a.Prev
	one_step.After = a.After
	return one_step
}

func (a *Oper) SetHook(prev, after chromeOper.Call) {
	a.Prev = prev
	a.After = after
}

func NewOper(name, oper, sel string, step int, args string, time_out int) *Oper {

	if oper == "" || step == 0 {
		return nil
	}

	o := new(Oper)
	o.Oper = oper
	o.Sel = sel
	o.Step = step
	o.Name = name
	o.Args = args
	o.TimeOut = time_out

	return o
}

//添加步骤
func (a *Opers) Add(name, oper, sel string, args string, step, timeOut int) error {
	index := len(*a)
	if step == 0 {
		if index == 0 {
			step = 5
		} else {
			step = (*a)[index-1].Step + 5
		}
	}

	o := NewOper(name, oper, sel, step, args, timeOut)
	if o != nil {
		*a = append(*a, o)
		return nil
	}
	return errors.New("添加失败")
}
