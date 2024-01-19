package config

import (
	"dog/pkg/validate"
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Address string `validate:"hostname_port"`

		// 值为 release 时：
		// 1. gin.Mode 也为 release
		// 2. swagger api 关闭
		Mode string `validate:"oneof=debug release"`
	}
	DB struct {
		DSN string
	}
	Logger struct {
		Level     string `validate:"oneof=debug info warn error"`
		Filename  string
		MaxSize   int
		MaxAge    int
		Compress  bool
		SplitTime int
	}
}

var instance = new(Config)

func Get() *Config {
	return instance
}

func Init(filename string) error {
	conf := viper.New()
	conf.SetConfigFile(filename)
	conf.SetConfigType("yaml")

	if err := conf.ReadInConfig(); err != nil {
		return fmt.Errorf("open config file err: %w", err)
	}

	// 加载到内存中
	if err := conf.Unmarshal(instance); err != nil {
		return fmt.Errorf("load config to memory err: %w", err)
	}

	// 校验配置
	if err := validate.Struct(instance); err != nil {
		return fmt.Errorf("validate config err: %w", err)
	}

	// 配置文件热更新
	conf.OnConfigChange(func(e fsnotify.Event) {
		if err := conf.Unmarshal(instance); err != nil {
			log.Printf("[dog] [ERROR] [config] file reload err: %s", err.Error())
		}
	})
	conf.WatchConfig()

	return nil
}
