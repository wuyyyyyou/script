package main

import (
	"github.com/sirupsen/logrus"

	"github.com/wuyyyyou/script/bai_sheng_cve/service"
)

func main() {
	svc := service.NewService()
	err := svc.TranslateDescription("/Users/leyouming/program/golang_note/golang_code/go_code/script/bai_sheng_cve/file/百胜-3.13-new-new-new.xlsx")
	if err != nil {
		logrus.Fatal(err)
	}
}
