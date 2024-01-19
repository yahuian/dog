## 功能

默认打印日志到终端，初始化时设置 FileOption 可以同时输出到终端和文件

```go
	// 增加文件相关配置
	file := &logger.FileOption{
		Filename:  "/var/log/app.log",
		MaxSize:   100,  // MB
		MaxAge:    30,   // days
		Compress:  true, // .gz
		SplitTime: 24,   // 每24h切割一次
	}

	logger.Init(file)
	defer logger.Sync()
```

## 分级规范

- Debug: 开发时使用，可能会打印详细的参数值，注意隐私信息不要打印（密码，api secret key 等）
- Info: 线上环境默认级别
- Warn: 一般是客户端请求错误时打印（如果频繁出现 warn 可能有人恶意攻击）
- Error: 一般是服务器内部错误时打印
