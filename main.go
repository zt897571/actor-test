// Package ztActor -----------------------------
// @file      : main.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/6/24 18:05
// -------------------------------------------
package main

import (
	"fmt"
	"time"
	"ztActor/engine/actor"
	_ "ztActor/engine/plugin"
	"ztActor/iface"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			fmt.Printf("err = %v \n", err)
		}
	}()
	err = iface.LoadPlugin("logic/plugin.0.so")
	if err != nil {
		return
	}
	userPid, err := actor.Spawn(iface.User)
	if err != nil {
		return
	}
	rst, err := actor.Call(userPid, 123, time.Second)
	if err != nil {
		return
	}
	fmt.Printf("rst1 = %v \n", rst)
	// 热更
	err = iface.LoadPlugin("logic/plugin.123.so")
	if err != nil {
		return
	}
	rst, err = actor.Call(userPid, 123, time.Second)
	if err != nil {
		return
	}
	fmt.Printf("rst2 = %v \n", rst)
}
