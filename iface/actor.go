// Package iface -----------------------------
// @file      : actor.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/8/2 10:50
// -------------------------------------------
package iface

import "time"

type IPid uint32
type ProcessStatus int
type ActorType int

type IProcessMgr interface {
	OnReload()
}

type IProcess interface {
	Id() IPid
	Status() ProcessStatus
	Call(target IPid, msg any, timeout time.Duration) (any, error)
	Cast(targetPid IPid, msg any) error
}

type IActor interface {
	SetProcess(process IProcess)
	GetData() any
	SetData(data any)
	Init(params ...interface{}) error
	OnStop(reason string) error
	HandleCall(sourcePid IPid, msg any) (any, error)
	HandleCast(sourcePid IPid, msg any) error
	GetActorType() ActorType
}

type IPluginMgr interface {
	LoadPlugin(pluginName string) error
}

const (
	Init ProcessStatus = iota
	Running
	Stoping
	Stoped
)

const (
	User ActorType = iota //
	Test
)
