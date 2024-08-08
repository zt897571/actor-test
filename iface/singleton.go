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

var actorMap = make(map[ActorType]func() IActor)

func RegisterActor(actorType ActorType, newActor func() IActor) {
	actorMap[actorType] = newActor
}

func GetActorByType(actorType ActorType) IActor {
	if newActor, ok := actorMap[actorType]; ok {
		return newActor()
	}
	return nil
}

func LoadPlugin(pluginName string) error {
	return G.PluginMgr.LoadPlugin(pluginName)
}
