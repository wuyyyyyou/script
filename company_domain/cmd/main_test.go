package cmd

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/wuyyyyou/script/company_domain/service"
)

// 更新数据库结构
func TestMain_AutoMigrate(t *testing.T) {
	svc := service.NewService()
	err := svc.DB.AutoMigrate(
		&service.DomainICPInfo{},
		//&service.CompanyInfo{},
	)
	if err != nil {
		logrus.Error(err)
	}
}

// 将文件中的域名写入数据库
func TestMain_WriteDomain(t *testing.T) {
	filePath := "/Users/leyouming/program/golang_note/golang_code/go_code/script/company_domain/file/清单2.xlsx"
	sheetName := "Sheet1"

	svc := service.NewService()
	err := svc.WriteDomainDao(filePath, sheetName)
	if err != nil {
		logrus.Error(err)
	}

}

// 测试通过站长之家的API获取网站备案信息
func TestMain_GetICPDomainInfoFromAPI(t *testing.T) {
	svc := service.NewService()

	companyInfo := new(service.CompanyInfo)
	err := svc.DB.Where("id = ?", "2").First(companyInfo).Error
	if err != nil {
		logrus.Error(err)
	}

	err = svc.GetICPDomainInfoFromAPI(companyInfo)
	if err != nil {
		logrus.Error(err)
	}
}

// 获取全部公司的域名备案信息
func TestMain_GetICPDomainInfo(t *testing.T) {
	svc := service.NewService()
	err := svc.GetICPDomainInfo()
	if err != nil {
		logrus.Error(err)
	}
}

func TestMain_ReadJsonAndWriteDomain(t *testing.T) {
	svc := service.NewService()
	err := svc.ReadJsonAndWriteDomain("上汽通用汽车有限公司")
	if err != nil {
		logrus.Error(err)
	}
}

func TestMain_ReadTxtAndWriteDomain(t *testing.T) {
	svc := service.NewService()
	err := svc.ReadTxtAndWriteDomain("斑马信息科技有限公司/斑马网络技术有限公司")
	if err != nil {
		logrus.Error(err)
	}
}

func TestMain_WriteExcel(t *testing.T) {
	svc := service.NewService()
	err := svc.WriteExcel("/Users/leyouming/program/golang_note/golang_code/go_code/script/company_domain/file/备案域名查询2.xlsx")
	if err != nil {
		logrus.Error(err)
	}
}
