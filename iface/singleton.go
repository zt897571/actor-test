// Package iface -----------------------------
// @file      : singleton.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/8/2 10:53
// -------------------------------------------
package iface

var G = &Singleton{}

type Singleton struct {
	ProcessMgr IProcessMgr // actor manager
	PluginMgr  IPluginMgr
}

func LoadPlugin(pluginName string) error {
	err := G.PluginMgr.LoadPlugin(pluginName)
	G.ProcessMgr.OnReload()
	return err
}

func RegisterActor(actorType ActorType, newActorFunc func() IActor) {
	G.ProcessMgr.RegisterActor(actorType, newActorFunc)
}
