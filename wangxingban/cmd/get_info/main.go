package main

import (
	"path/filepath"

	service "github.com/wuyyyyou/script/wangxingban"

	"github.com/sirupsen/logrus"
)

func main() {
	service.Logger.SetLevel(logrus.DebugLevel)

	fileRootDir := "/Users/leyouming/program/golang_note/golang_code/go_code/script/wangxingban/files"

	srcExcelPath := filepath.Join(fileRootDir, "20240322", "上海大岂网络科技有限公司.xlsx")
	dstExcelPath := filepath.Join(fileRootDir, "20240322", "上海大岂网络科技有限公司-更新.xlsx")

	svc := service.NewService(srcExcelPath, dstExcelPath)
	defer svc.Close()
	err := svc.Worker()
	if err != nil {
		service.Logger.Error(err)
	}
}
