package config

import (
	"fmt"
	"github.com/spf13/viper"
	"http2db/utils"
	"path"
	"runtime"
	"sync"
)

// 自定义 viper 读取数据变量用 mapstructuretag tag
type Config struct {
	DBConfig           *DBConfig
	MaxPool            int
	SpecialChar        bool
	QuickFilteringMode bool
	DataAnalyze        struct {
		MaxSemaphore        int
		MapperMaxGoroutines int
	}
}

var instance *Config
var once sync.Once

func Global() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

// 隐式初始化 Config 文件
func init() {

	// 获取当前文件的路径
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Dir(path.Dir(filename))
	appDir := utils.GetAppPath()

	viper.SetConfigName("tool") // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	//viper.AddConfigPath("/etc/dsu")  // path to look for the config file in
	//viper.AddConfigPath("$HOME/dsu") // call multiple times to add many search paths
	viper.AddConfigPath(appDir)
	viper.AddConfigPath(rootDir) // optionally look for config in the working directory
	err := viper.ReadInConfig()  // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// 绑定环境变量
	viper.AutomaticEnv()

	// 绑定全局变量
	err = viper.Unmarshal(Global())
	if err != nil {
		panic(fmt.Sprintf("unable to decode into struct, %v", err))
	}
}
