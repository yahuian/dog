## 状态码 & 业务码

采用常见的标准 HTTP Status Code + HTTP Response Body 中业务自定义 code

分为以下三大类，200、400、500

- 200 正常
- 400 客户端请求错误
  - 401 未登录
  - 403 没有权限
  - 404 资源未找到
  - 40001 - 49999 业务细分错误
    - 41xxx 代表 x 模块
    - 42xxx 代表 y 模块
    - ...
- 500 服务端内部错误

```json
{
  "Code": 41001,
  "Msg": "密码强度不符合要求",
}
```

## 开发规范

### 所有服务端的报错信息不可直接暴露出去

因为某些错误信息可能包含底层数据库的信息，文件路径等等，防止有人猜测系统内部结构，并以此来攻击系统

### 错误处理采用人肉增加上下文形成错误链

`fmt.Errorf("some context msg: %w", err)`

最终的效果直观清晰，不看代码也能知道问题出在哪里了

```
create user err: mysql connect err: dial 192.168.1.100 err: ...
```

（后续等等这个：https://zhuanlan.zhihu.com/p/670897856）

### 错误第一次出现的地方必须增加 errcode 区分清楚类型

比如客户端请求出错 `return errcode.BadRequest(err)`

服务器内部错误 `return errcode.Server(err)`

业务自定义错误 `return errcode.UserStatusBannedErr`

（api response 会根据其中的 code 来区分错误类型）

### 新增的自定义错误必须在本模块下

```go
UserStatusBannedErr = newErrcode(41000, "账号已被封禁")
...
```

### error 只在顶层处理（如 api 层），service 和 model 层直接返回即可

```go
	// service 层

	total, err := userModel.Count(ctx, query, args)
	if err != nil {
		log.Println(err) // 不要重复打印
		return err
	}
```

### 程序启动时发生严重错误（数据库无法连接，配置加载失败等）直接 panic，方便立马排查错误

```go
	if err := config.Init(configFile); err != nil {
		panic(err)
	}
```
