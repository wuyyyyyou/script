package service

import (
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/remeh/sizedwaitgroup"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	ServiceConfig *ServiceConfig `toml:"service"`
	MysqlConfig   *MysqlConfig   `toml:"mysql"`
}

type ServiceConfig struct {
	IsDebug           bool   `toml:"is_debug"`
	HttpMaxGoroutines int    `toml:"http_max_goroutines"`
	IPInfoKey         string `toml:"ip_info_key"`
}

type MysqlConfig struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Dbname   string `toml:"dbname"`
	Timeout  string `toml:"timeout"`
}

type Service struct {
	Cache             sync.Map
	DB                *gorm.DB
	HttpMaxGoroutines int
	swg               sizedwaitgroup.SizedWaitGroup

	IPInfoKey string
}

func NewService() *Service {
	config := new(Config)
	_, err := toml.DecodeFile("config.toml", config)
	if err != nil {
		panic(fmt.Sprintf("读取toml配置文件错误:%s", err))
	}

	// 日志水平
	if config.ServiceConfig.IsDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	service := new(Service)

	// 设置http访问最大并发数
	service.HttpMaxGoroutines = config.ServiceConfig.HttpMaxGoroutines
	service.swg = sizedwaitgroup.New(service.HttpMaxGoroutines)

	service.DB = getDB(config.MysqlConfig)

	service.IPInfoKey = config.ServiceConfig.IPInfoKey

	return service
}

func getDB(mysqlConfig *MysqlConfig) *gorm.DB {
	username := mysqlConfig.Username
	password := mysqlConfig.Password
	host := mysqlConfig.Host
	port := mysqlConfig.Port
	Dbname := mysqlConfig.Dbname
	timeout := mysqlConfig.Timeout

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s",
		username, password, host, port, Dbname, timeout)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	return db
}
