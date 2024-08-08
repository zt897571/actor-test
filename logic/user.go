// Package main -----------------------------
// @file      : User.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/8/2 19:07
// -------------------------------------------
package main

import (
	"fmt"
	"ztActor/iface"
	"ztActor/logic_data"
)

var _ iface.IActor = (*User)(nil)

func init() {
	fmt.Printf("plugin function\n")
	iface.RegisterActor(iface.User, newUser)
}

type User struct {
	iface.IProcess
	*logic_data.UserData
}

func newUser() iface.IActor {
	return &User{UserData: &logic_data.UserData{}}
}

func (u *User) SetProcess(process iface.IProcess) {
	u.IProcess = process
}

func (u *User) GetData() any {
	return u.UserData
}

func (u *User) SetData(data any) {
	u.UserData = data.(*logic_data.UserData)
}

func (u *User) Init(params ...interface{}) error {
	return nil
}

func (u *User) OnStop(reason string) error {
	fmt.Printf("User stoped reason = %s \n", reason)
	return nil
}

func (u *User) HandleCall(sourcePid iface.IPid, msg any) (any, error) {
	//u.age = msg
	fmt.Printf("test = %d\n", 3333)
	return 33333, nil
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
