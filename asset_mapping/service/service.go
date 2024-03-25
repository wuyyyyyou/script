package service

import (
	"fmt"
	"sync"

	"github.com/remeh/sizedwaitgroup"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Service struct {
	Cache sync.Map

	DB *gorm.DB

	HttpMaxGoroutines int
	swg               sizedwaitgroup.SizedWaitGroup
}

func NewService() *Service {
	// 日志水平
	logrus.SetLevel(logrus.DebugLevel)

	service := new(Service)

	// 设置http访问最大并发数
	service.HttpMaxGoroutines = 10
	service.swg = sizedwaitgroup.New(service.HttpMaxGoroutines)

	service.DB = getDB()

	return service
}

func getDB() *gorm.DB {
	username := "root"
	password := "123456"
	host := "127.0.0.1"
	port := 3306
	Dbname := "asset_mapping"
	timeout := "10s"

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
