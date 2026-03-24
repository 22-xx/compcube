package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Configuration struct {
	ConfigViper *viper.Viper

	App    App      `mapstructure:"app" json:"app" yaml:"app"`
	DB     Database `mapstructure:"database" json:"database" yaml:"database"`
	Path   Path     `mapstructure:"file_path" json:"file_path" yaml:"file_path"`
	Docker Docker   `mapstructure:"docker_path" json:"docker_path" yaml:"docker_path"`
}

type App struct {
	Env     string `mapstructure:"env" json:"env" yaml:"env"`
	Host    string `mapstructure:"host" json:"host" yaml:"host"`
	Port    string `mapstructure:"port" json:"port" yaml:"port"`
	AppName string `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
	AppUrl  string `mapstructure:"app_url" json:"app_url" yaml:"app_url"`
}

type Database struct {
	Name     string `mapstructure:"name" json:"name" yaml:"name"`
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

type Path struct {
	RootPath string `mapstructure:"root_path" json:"root_path" yaml:"root_path"`
	LogPath  string `mapstructure:"log_path" json:"log_path" yaml:"log_path"`
	LogName  string `mapstructure:"log_name" json:"log_name" yaml:"log_name"`
	FilePath string `mapstructure:"file_path" json:"file_path" yaml:"file_path"`
}

type Docker struct {
	InputPath  string `mapstructure:"input_path" json:"input_path" yaml:"input_path"`
	UserPath   string `mapstructure:"user_path" json:"user_path" yaml:"user_path"`
	OutputPath string `mapstructure:"output_path" json:"output_path" yaml:"output_path"`
	LogPath    string `mapstructure:"log_path" json:"log_path" yaml:"log_path"`
}

var Config Configuration

func ConfigInit() {
	v := viper.New()
	v.SetConfigFile("./config/config.yaml")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %s \n", err))
	}

	// 监听配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
		// 重载配置
		if err := v.Unmarshal(&Config); err != nil {
			fmt.Println(err)
		}
	})
	// 将配置赋值给全局变量
	if err := v.Unmarshal(&Config); err != nil {
		fmt.Println(err)
	}
}
