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

var _ iface.IProcessMgr = (*processMgr)(nil)

type processMgr struct {
	idx                uint32
	processMap         sync.Map
	actorCreateFuncMap sync.Map
}

func newProcessMgr() *processMgr {
	return &processMgr{
		idx:        0,
		processMap: sync.Map{},
	}
}

func (m *processMgr) spawn(actor iface.IActor, params ...interface{}) (iface.IPid, error) {
	id := atomic.AddUint32(&m.idx, 1)
	pid := iface.IPid(id)
	p := newProcess(pid, actor)
	m.processMap.Store(id, p)
	return pid, p.Start(params)
}

func (m *processMgr) stop(pid iface.IPid, reason string, block bool) error {
	if p, ok := m.processMap.Load(uint32(pid)); ok {
		return p.(*process).Stop(reason, block)
	} else {
		return engine.ErrProcessNotFound
	}
}

func (m *processMgr) getProcess(pid iface.IPid) iface.IProcess {
	if pr, isOk := m.processMap.Load(uint32(pid)); isOk {
		return pr.(iface.IProcess)
	}
	return nil
}

func (m *processMgr) OnReload() {
	m.processMap.Range(func(key, value interface{}) bool {
		err := value.(*process).reload()
		if err != nil {
			fmt.Printf("reload process error %v \n", err)
		}
		return true
	})
}

func (m *processMgr) RegisterActor(actorType iface.ActorType, newActor func() iface.IActor) {
	m.actorCreateFuncMap.Store(actorType, newActor)
}

func (m *processMgr) CreateActor(actorType iface.ActorType) iface.IActor {
	if newActor, ok := m.actorCreateFuncMap.Load(actorType); ok {
		return newActor.(func() iface.IActor)()
	}
	return nil
}

func routeMsg(target iface.IPid, msg *ProcessMsg) error {
	if p := iface.G.ProcessMgr.(*processMgr).getProcess(target); p != nil {
		return p.(*process).pushMsg(msg)
	} else {
		//todo  接入rpc
		// remote  ?
	}
	return nil
}

func Spawn(actType iface.ActorType, params ...interface{}) (iface.IPid, error) {
	mgr := iface.G.ProcessMgr.(*processMgr)
	actor := mgr.CreateActor(actType)
	if actor == nil {
		return 0, engine.ErrActorNotFound
	}
	return mgr.spawn(actor, params)
}

func Stop(pid iface.IPid, reason string, block bool) error {
	return getMgr().stop(pid, reason, block)
}

func getMgr() *processMgr {
	return iface.G.ProcessMgr.(*processMgr)
}
