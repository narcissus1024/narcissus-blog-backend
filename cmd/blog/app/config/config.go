package config

import (
	"path/filepath"

	"github.com/narcissus1949/narcissus-blog/internal/database/cache"
	"github.com/narcissus1949/narcissus-blog/internal/database/mysql"
	"github.com/narcissus1949/narcissus-blog/internal/logger"
	"github.com/narcissus1949/narcissus-blog/internal/utils"
	"github.com/spf13/viper"
)

var (
	viperConfig = viper.New()
	Config      = NewConfig()
)

type Conf struct {
	App    AppConfig         `json:"app"`
	Mysql  mysql.MysqlConfig `json:"mysql"`
	Redis  cache.RedisConfig `json:"redis"`
	Logger logger.LogConfig  `json:"logger"`
}

type AppConfig struct {
	Name          string `json:"name"`
	Port          int    `json:"port"`
	ImgDataDir    string `json:"imgDataDir"`
	ImgProxyURL   string `json:"imgProxyURL"`
	PrivateKeyDir string `json:"privateKeyDir"`
	PublicKeyDir  string `json:"publicKeyDir"`
}

func NewDefaultAppCfg() AppConfig {
	rootDir := utils.GetRootDir()
	return AppConfig{
		Name:          "blog",
		Port:          8080,
		ImgDataDir:    filepath.Join(rootDir, "data", "img"),
		ImgProxyURL:   "http://127.0.0.1:8082/img",
		PrivateKeyDir: filepath.Join(rootDir, "data", "conf"),
		PublicKeyDir:  filepath.Join(rootDir, "data", "conf"),
	}
}

func NewConfig() Conf {
	return Conf{
		App:    NewDefaultAppCfg(),
		Logger: logger.NewDefaultCfg(),
		Mysql:  mysql.NewDefaultMysqlCfg(),
		Redis:  cache.NewDefaultRedisCfg(),
	}
}

func MustInit(fliePath string) {
	viperConfig.SetConfigFile(fliePath)
	if err := viperConfig.ReadInConfig(); err != nil {
		panic(err.Error())
	}

	if err := viperConfig.Unmarshal(&Config); err != nil {
		panic(err.Error())
	}
	if err := Config.Check(); err != nil {
		panic(err.Error())
	}
}

// todo
func (c *Conf) Check() error {
	return nil
}
