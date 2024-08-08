// Package ztActor -----------------------------
// @file      : main.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/6/24 18:05
// -------------------------------------------
package main

import (
	"bufio"
	"fmt"
	"os"
	"plugin"
	"regexp"
	"time"
	"ztActor/engine/actor"
	_ "ztActor/engine/plugin"
	"ztActor/iface"
)

var userPid iface.IPid
var reloadRe, _ = regexp.Compile(`reload\s+(\d+)`)

func main() {
	_, err := plugin.Open("logic/plugin.0.so")
	if err != nil {
		return
	}
	userPid, err = actor.Spawn(iface.User)
	if err != nil {
		return
	}
	fmt.Printf("pid = %d \n", userPid)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("> ")
		reRst := reloadRe.FindStringSubmatch(line)
		if len(reRst) > 0 {
			version := reRst[1]
			pluginName := fmt.Sprintf("logic/plugin.%s.so", version)
			fmt.Printf("reload pluginName = %s \n", pluginName)
			err = iface.LoadPlugin(pluginName)
			if err != nil {
				fmt.Printf("reloaded error = %s \n", err)
				return
			}
			// todo zhangtuo 改为监听者模式
			iface.G.ProcessMgr.OnReload()
			fmt.Println("reloaded")
		} else if line == "call" {
			rst, err := actor.Call(userPid, 123, time.Second)
			fmt.Printf("rst = %v err = %v \n", rst, err)
		}
	}
}
