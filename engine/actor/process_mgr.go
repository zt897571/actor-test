// Package actor -----------------------------
// @file      : process_mgr.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/8/2 10:51
// -------------------------------------------
package actor

import (
	"fmt"
	"sync"
	"sync/atomic"
	"ztActor/engine"
	"ztActor/iface"
)

func init() {
	iface.G.ProcessMgr = newProcessMgr()
}

var _ iface.IProcessMgr = (*ProcessMgr)(nil)

type ProcessMgr struct {
	idx        uint32
	processMap sync.Map
}

func newProcessMgr() *ProcessMgr {
	return &ProcessMgr{
		idx:        0,
		processMap: sync.Map{},
	}
}

func (m *ProcessMgr) spawn(actor iface.IActor, params ...interface{}) (iface.IPid, error) {
	id := atomic.AddUint32(&m.idx, 1)
	pid := iface.IPid(id)
	p := newProcess(pid, actor)
	m.processMap.Store(id, p)
	return pid, p.Start(params)
}

func (m *ProcessMgr) stop(pid iface.IPid, reason string, block bool) error {
	if p, ok := m.processMap.Load(uint32(pid)); ok {
		return p.(*process).Stop(reason, block)
	} else {
		return engine.ErrProcessNotFound
	}
}

func (m *ProcessMgr) getProcess(pid iface.IPid) iface.IProcess {
	if pr, isOk := m.processMap.Load(uint32(pid)); isOk {
		return pr.(iface.IProcess)
	}
	return nil
}

func (m *ProcessMgr) OnReload() {
	m.processMap.Range(func(key, value interface{}) bool {
		err := value.(*process).reload()
		if err != nil {
			fmt.Printf("reload process error %v \n", err)
		}
		return true
	})
}

func routeMsg(target iface.IPid, msg *ProcessMsg) error {
	if p := iface.G.ProcessMgr.(*ProcessMgr).getProcess(target); p != nil {
		return p.(*process).pushMsg(msg)
	} else {
		// remote  ?
	}
	return nil
}

func Spawn(actType iface.ActorType, params ...interface{}) (iface.IPid, error) {
	actor := iface.GetActorByType(actType)
	if actor == nil {
		return 0, engine.ErrUnknown
	}

	return iface.G.ProcessMgr.(*ProcessMgr).spawn(actor, params)
}

func Stop(pid iface.IPid, reason string, block bool) error {
	return iface.G.ProcessMgr.(*ProcessMgr).stop(pid, reason, block)
}
