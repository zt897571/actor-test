// Package plugin -----------------------------
// @file      : plugin.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/8/8 16:58
// -------------------------------------------
package plugin

import (
	"fmt"
	"plugin"
	"regexp"
	"strconv"
	"sync"
	"ztActor/engine"
	"ztActor/iface"
)

func init() {
	iface.G.PluginMgr = newMgr()
}

var _ iface.IPluginMgr = (*mgr)(nil)

type mgr struct {
	mutext         sync.Mutex
	currendVersion int32
}

func newMgr() iface.IPluginMgr {
	return &mgr{
		currendVersion: -1,
	}
}

func (m *mgr) LoadPlugin(filePath string) error {
	m.mutext.Lock()
	defer m.mutext.Unlock()

	pluginName, version, err := GetPluginVersion(filePath)
	if err != nil {
		return err
	}
	if version <= m.currendVersion {
		return engine.ErrPluginVersionIsSmaller
	}
	fmt.Printf("load plugin %s version = %d\n", pluginName, version)
	_, err = plugin.Open(filePath)
	if err != nil {
		return err
	}
	m.currendVersion = version
	return nil
}

var rg, _ = regexp.Compile(`plugin\.(\d+)\.so`)

// GetPluginVersion 获取plugin 版本号
// version ->   plugin.1232.so
func GetPluginVersion(pluginName string) (string, int32, error) {
	match := rg.FindStringSubmatch(pluginName)
	if match != nil && len(match) == 2 {
		version, err := strconv.ParseInt(match[1], 10, 32)
		if err != nil {
			return "", 0, err
		}
		return match[0], int32(version), nil
	}
	return "", 0, engine.ErrPluginNameFormatError
}
