// Package plugin -----------------------------
// @file      : plugin.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/8/8 16:58
// -------------------------------------------
package plugin

import (
	"plugin"
	"regexp"
	"strconv"
	"ztActor/engine"
	"ztActor/iface"
)

func init() {
	iface.G.PluginMgr = newMgr()
}

type Mgr struct {
	reloaded       map[string]struct{}
	currendVersion uint32
	isStart        bool
}

func newMgr() iface.IPluginMgr {
	return &Mgr{
		reloaded: map[string]struct{}{},
	}
}

func (m *Mgr) LoadPlugin(pluginName string) error {
	if _, ok := m.reloaded[pluginName]; ok {
		return engine.ErrPluginAlreadyLoaded
	}
	version, err := GetPluginVersion(pluginName)
	if err != nil {
		return err
	}
	if version <= m.currendVersion {
		return engine.ErrPluginVersionIsSmaller
	}
	_, err = plugin.Open(pluginName)
	if err != nil {
		return err
	}
	m.reloaded[pluginName] = struct{}{}
	m.currendVersion = version
	return nil
}

var rg, _ = regexp.Compile(`plugin\.(\d+)\.so`)

// GetPluginVersion 获取plugin 版本号
// version ->   plugin.1232.so
func GetPluginVersion(pluginName string) (uint32, error) {
	match := rg.FindStringSubmatch(pluginName)
	if match != nil {
		version, err := strconv.ParseUint(match[1], 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(version), nil
	}
	return 0, engine.ErrPluginNameFormatError
}
