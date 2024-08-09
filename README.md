#zt test


actor 模型 + 热更新

## 示例
main.go 函数中已经有了基础示例

测试步骤：
1. 调用make build会生成test 和plugin.0.so两个文件
2. 修改 user.go 中的handle call 返回值例如12312312, 执行make version=123 buildPlugin生成 plugin.123.so
3. 执行main函数会显示热更前和热更后的结果


## 说明

1. actor相关
暂未实现call其他节点的actor操作，只实现了本地actor的调用, 后续接入rpc框架实现远程actor调用

示例：
```golang
// 创建actor
pid, err := actor.Spawn(iface.User)

// actor call
actor.Call(pid, 123)

// actor cast
actor.Cast(pid, 123)
```

2. 热更新相关
 借助golang的plugin, 模仿erlang的热更新机制，通过数据和函数分离，动态替换actor回调类实现代码热更新







