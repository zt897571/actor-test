// Package actor -----------------------------
// @file      : actor.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/8/2 11:00
// -------------------------------------------
package actor

import (
	"fmt"
	"sync/atomic"
	"time"
	"ztActor/engine"
	"ztActor/iface"
)

type processMsgType int

var _ iface.IProcess = (*process)(nil)

const (
	StopMsg = iota
	CallMsg // call 请求
	CastMsg // cast 请求
	ReloadMsg
)

type ProcessMsg struct {
	msgType processMsgType
	param   interface{}
}

type process struct {
	iface.IActor
	id      iface.IPid
	status  int32
	mailBox chan *ProcessMsg
}

func newProcess(id iface.IPid, actor iface.IActor) *process {
	p := &process{
		id:      id,
		IActor:  actor,
		mailBox: make(chan *ProcessMsg, 10000),
	}
	p.setStatus(iface.Init)
	actor.SetProcess(p)
	return p
}

func (p *process) Id() iface.IPid {
	return p.id
}

func (p *process) Start(args ...interface{}) error {
	if p.Status() != iface.Init {
		return engine.ErrProcessStatusError
	}
	startChan := make(chan error)
	go func() {
		if p.Status() != iface.Init {
			startChan <- engine.ErrProcessStatusError
			return
		}
		err := p.Init(args)
		if err != nil {
			startChan <- err
			return
		}
		p.setStatus(iface.Running)
		startChan <- nil
		for {
			select {
			case msg := <-p.mailBox:
				err = p.handleMsg(msg)
				if err != nil {
					fmt.Printf("handle msg error = %s", err)
				} else if msg.msgType == StopMsg {
					return
				}
			}
		}
	}()
	return <-startChan
}

func (p *process) pushMsg(msg *ProcessMsg) error {
	if p.Status() != iface.Running {
		return engine.ErrProcessStatusError
	}
	p.mailBox <- msg
	return nil
}

func (p *process) Stop(reason string, block bool) error {
	stMsg := &stopMsg{reason: reason}
	msg := &ProcessMsg{
		msgType: StopMsg,
		param:   stMsg,
	}
	if block {
		stMsg.stopChan = make(chan error)
		err := p.pushMsg(msg)
		if err != nil {
			return err
		}
		return <-stMsg.stopChan
	} else {
		return p.pushMsg(msg)
	}
}

func (p *process) setStatus(status iface.ProcessStatus) {
	atomic.StoreInt32(&p.status, int32(status))
}

func (p *process) Status() iface.ProcessStatus {
	return iface.ProcessStatus(atomic.LoadInt32(&p.status))
}

// 消息处理
func (p *process) handleMsg(msg *ProcessMsg) error {
	switch msg.msgType {
	case StopMsg:
		if p.Status() != iface.Running {
			return engine.ErrProcessStatusError
		}
		stMsg := msg.param.(*stopMsg)
		p.setStatus(iface.Stoping)
		err := p.OnStop(stMsg.reason)
		p.setStatus(iface.Stoped)
		if stMsg.stopChan != nil {
			stMsg.stopChan <- err
		}
		return err
	case CallMsg:
		callMsg := msg.param.(*callMsg)
		reply, err := p.HandleCall(callMsg.source, callMsg.msg)
		callMsg.waitChan <- reply
		return err
	case CastMsg:
		castMsg := msg.param.(*castMsg)
		return p.HandleCast(castMsg.source, castMsg.msg)
	case ReloadMsg:
		actor := getMgr().CreateActor(p.GetActorType())
		if actor == nil {
			return engine.ErrActorNotFound
		}
		data := p.IActor.GetData()
		actor.SetData(data)
		p.IActor = actor
		fmt.Println("reload succ")
		return nil
	}
	return nil
}

func (p *process) Call(targetPid iface.IPid, req any, timeout time.Duration) (any, error) {
	if targetPid == p.Id() {
		return nil, engine.ErrCanNotCallSelf
	}
	return internalCall(p.Id(), targetPid, req, timeout)
}

func (p *process) Cast(targetPid iface.IPid, req any) error {
	msg := &castMsg{source: p.Id(), target: targetPid, msg: req}
	return routeMsg(targetPid, &ProcessMsg{msgType: CastMsg, param: msg})
}

func (p *process) reload() error {
	return p.pushMsg(&ProcessMsg{msgType: ReloadMsg})
}

func Call(targetPid iface.IPid, req any, timeout time.Duration) (interface{}, error) {
	return internalCall(0, targetPid, req, timeout)
}

func Cast(targetPid iface.IPid, req any) error {
	return internalCast(0, targetPid, req)
}

func internalCast(sourcePid, targetPid iface.IPid, msg any) error {
	castMsg := &castMsg{source: sourcePid, target: targetPid, msg: msg}
	return routeMsg(targetPid, &ProcessMsg{msgType: CastMsg, param: castMsg})
}

func internalCall(sourcePid, targetPid iface.IPid, msg any, timeout time.Duration) (any, error) {
	waitChan := make(chan any)
	callmsg := &callMsg{source: sourcePid, target: targetPid, waitChan: waitChan, msg: msg}
	err := routeMsg(targetPid, &ProcessMsg{msgType: CallMsg, param: callmsg})
	if err != nil {
		return nil, err
	}
	select {
	case reply := <-waitChan:
		return reply, nil
	case <-time.After(timeout):
		return nil, engine.ErrTimeout
	}
}
