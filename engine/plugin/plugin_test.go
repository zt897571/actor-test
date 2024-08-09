// Package plugin -----------------------------
// @file      : plugin_test.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/8/8 17:20
// -------------------------------------------
package plugin

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPluginName(t *testing.T) {
	convey.Convey("test plugin Name", t, func() {
		name, version, err := GetPluginVersion("logic/plugin.1232.so")
		convey.So(name, convey.ShouldEqual, "plugin.1232.so")
		convey.So(err, convey.ShouldEqual, nil)
		convey.So(version, convey.ShouldEqual, 1232)
	})
}
