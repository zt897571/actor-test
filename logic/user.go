// Package main -----------------------------
// @file      : User.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/8/2 19:07
// -------------------------------------------
package main

import (
	"fmt"
	"ztActor/data"
	"ztActor/iface"
)

var _ iface.IActor = (*User)(nil)

func init() {
	fmt.Printf("plugin init\n")
	iface.RegisterActor(iface.User, newUser)
}

type User struct {
	iface.IProcess
	*data.UserData
}

func newUser() iface.IActor {
	return &User{UserData: &data.UserData{}}
}

func (u *User) SetProcess(process iface.IProcess) {
	u.IProcess = process
}

func (u *User) GetData() any {
	return u.UserData
}

func (u *User) SetData(uData any) {
	u.UserData = uData.(*data.UserData)
}

func (u *User) Init(params ...interface{}) error {
	return nil
}

func (u *User) OnStop(reason string) error {
	fmt.Printf("User stoped reason = %s \n", reason)
	return nil
}

func (u *User) HandleCall(sourcePid iface.IPid, msg any) (any, error) {
	return 3333, nil
}

func (u *User) HandleCast(sourcePid iface.IPid, msg any) error {
	u.Age = msg.(int)
	return nil
}

func (u *User) GetActorType() iface.ActorType {
	return iface.User
}

func main() {
	fmt.Printf("plugin main\n")
}
