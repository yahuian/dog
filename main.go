package main

import (
	"dog/api"
	"dog/config"
	"dog/model"
	"dog/pkg/logger"
	"dog/pkg/validate"
	"flag"

	"github.com/gin-gonic/gin"
)

func main() {
	// 系统初始化
	mustInit()

	defer func() {
		logger.Sync()
	}()

	gin.SetMode(config.Get().Server.Mode)
	engine := gin.Default()

	// 注册路由
	api.Register(engine)

	// HTTP server
	if err := engine.Run(config.Get().Server.Address); err != nil {
		panic(err)
	}
}

// 需要注意各组件的初始化顺序
func mustInit() {
	// 参数校验
	if err := validate.Init(); err != nil {
		panic(err)
	}

	// 配置文件
	var filename string
	flag.StringVar(&filename, "config", "config/config.yaml", "config filename")
	flag.Parse()
	if err := config.Init(filename); err != nil {
		panic(err)
	}

	// 连接数据库
	if err := model.Init(); err != nil {
		panic(err)
	}

	// 日志
	conf := config.Get().Logger
	opt := &logger.FileOption{
		Filename:  conf.Filename,
		MaxSize:   conf.MaxSize,
		MaxAge:    conf.MaxAge,
		Compress:  conf.Compress,
		SplitTime: conf.SplitTime,
	}
	logger.Init(opt)
	logger.SetLevel(conf.Level)
}
