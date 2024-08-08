// Package actor -----------------------------
// @file      : msg.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/8/2 16:53
// -------------------------------------------
package actor

import "ztActor/iface"

type stopMsg struct {
	reason   string
	stopChan chan error
}

type callMsg struct {
	source   iface.IPid
	target   iface.IPid
	waitChan chan any
	msg      interface{}
}

type castMsg struct {
	source iface.IPid
	target iface.IPid
	msg    interface{}
}
