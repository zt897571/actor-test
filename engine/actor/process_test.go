// Package actor -----------------------------
// @file      : process_test.go
// @author    : zhangtuo
// @contact   :
// @time      : 2024/8/2 15:36
// -------------------------------------------
package actor

import (
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
	"ztActor/engine"
	"ztActor/iface"
)

var _ iface.IActor = (*testActor)(nil)

type testMsg struct {
	name   string
	param  any
	param1 any
}

type testActor struct {
	iface.IProcess
	flag uint32
}

func (t *testActor) GetActorType() iface.ActorType {
	return iface.Test
}

func (t *testActor) GetData() any {
	return t.flag
}

func (t *testActor) SetData(data any) {
	t.flag = data.(uint32)
}

func (t *testActor) HandleCall(sourcePid iface.IPid, msg any) (any, error) {
	tmsg := msg.(*testMsg)
	switch tmsg.name {
	case "echo":
		return tmsg.param, nil
	case "sleep":
		time.Sleep(tmsg.param.(time.Duration))
		return nil, nil
	case "call":
		pid := tmsg.param.(iface.IPid)
		msg := tmsg.param1
		reply, err := t.Call(pid, &testMsg{name: "echo", param: msg}, time.Second)
		return reply, err
	default:
		return nil, engine.ErrUnknown
	}
}

func (t *testActor) HandleCast(sourcePid iface.IPid, msg any) error {
	tmsg := msg.(*testMsg)
	switch tmsg.name {
	case "cast_echo":
		tmsg.param1.(chan int32) <- tmsg.param.(int32)
		return nil
	default:
		return engine.ErrUnknown
	}
}

func (t *testActor) Init(params ...interface{}) error {
	fmt.Println("start")
	time.Sleep(time.Second)
	return nil
}

func (t *testActor) OnStop(reason string) error {
	fmt.Println("stoping")
	time.Sleep(time.Second)
	return nil
}

func (t *testActor) SetProcess(process iface.IProcess) {
	t.IProcess = process
}

func TestActorStop(t *testing.T) {
	convey.Convey("test actor stop", t, func() {
		pr := newProcess(1, &testActor{})
		err := pr.Start()
		convey.So(err, convey.ShouldBeNil)
		err = pr.Stop("test", true)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestProcessManager_Spawn(t *testing.T) {
	convey.Convey("test process manager spawn", t, func() {
		pid, err := Spawn(iface.Test, 123)
		convey.So(err, convey.ShouldBeNil)
		err = Stop(pid, "test", true)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestProcess_Call(t *testing.T) {
	convey.Convey("test process call", t, func() {
		pid1, err := Spawn(iface.Test)
		convey.So(err, convey.ShouldBeNil)
		pid2, err := Spawn(iface.Test)
		convey.So(err, convey.ShouldBeNil)
		// test echo
		req := 123
		reply, err := Call(pid1, &testMsg{name: "echo", param: req}, time.Second)
		convey.So(err, convey.ShouldBeNil)
		convey.So(reply, convey.ShouldEqual, req)

		// test seq call
		reply, err = Call(pid1, &testMsg{name: "call", param: pid2, param1: req}, time.Second)
		convey.So(err, convey.ShouldBeNil)
		convey.So(reply, convey.ShouldEqual, req)

		// test sleep
		reply, err = Call(pid1, &testMsg{name: "sleep", param: time.Second * 5}, time.Second)
		convey.So(err, convey.ShouldEqual, engine.ErrTimeout)
		convey.So(reply, convey.ShouldEqual, nil)

	})
}

func TestProcess_Cast(t *testing.T) {
	convey.Convey("test process cast", t, func() {
		pid1, err := Spawn(iface.Test, 123)
		convey.So(err, convey.ShouldBeNil)
		var req int32 = 123
		resultChan := make(chan int32)
		err = Cast(pid1, &testMsg{name: "cast_echo", param: req, param1: resultChan})
		convey.So(err, convey.ShouldBeNil)
		convey.So(<-resultChan, convey.ShouldEqual, req)
	})
}
