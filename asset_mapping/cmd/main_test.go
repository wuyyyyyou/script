package cmd

import (
	"fmt"
	"testing"

	"github.com/wuyyyyou/script/asset_mapping/models"
	"github.com/wuyyyyou/script/asset_mapping/service"
)

// 数据库初始化
func TestMain_AutoMigrate(t *testing.T) {
	svc := service.NewService()
	err := svc.DB.AutoMigrate(
		//&models.Seed{},
		&models.IP{},
		//&models.Port{},
		//&models.Domain{},
		//&models.Company{},
	)
	if err != nil {
		t.Fatal(err)
	}
}

// 从excel文件中读取种子url，并存入数据库
func TestMain_ReadSeed(t *testing.T) {
	svc := service.NewService()
	err := svc.ReadSeed(
		"/Users/leyouming/company_program/script/asset_mapping/file/测绘2.xlsx",
		"IP",
		"IP-1", "IP-2")
	if err != nil {
		t.Fatal(err)
	}
}

// 读取数据库中的种子url生成Quake的查询语句
func TestMain_GenerateQuakeQuery(t *testing.T) {
	svc := service.NewService()
	query, err := svc.GenerateQuakeQuery(true)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("'%s'", query)
}

// 读取nmap的结果xml文件并将结果IP和Port存入数据库
func TestMain_ReadNmapXml(t *testing.T) {
	svc := service.NewService()
	err := svc.ReadNmapXml("/Users/leyouming/company_program/script/asset_mapping/file/test1.xml")
	if err != nil {
		t.Fatal(err)
	}
}

// 读取某个目录下的所有nmap的结果xml文件并将IP和Port存入数据库
func TestMain_ReadAllNmapXml(t *testing.T) {
	svc := service.NewService()
	err := svc.ReadAllNmapXml("/Users/leyouming/company_program/scan_tool/AssetMapping/nmap/0325_2")
	if err != nil {
		t.Fatal(err)
	}
}

// 读取quake导出的结果excel将子域名存储到数据库
func TestMain_ReadSubDomain(t *testing.T) {
	svc := service.NewService()
	err := svc.ReadSubDomain("/Users/leyouming/company_program/script/asset_mapping/file/服务数据_20240324_191828.xlsx")
	if err != nil {
		t.Fatal(err)
	}
}

// 获取数据库中所有子域名的IP
func TestMain_GetIP(t *testing.T) {
	svc := service.NewService()
	err := svc.GetIP(true)
	if err != nil {
		t.Fatal(err)
	}
}

// nmap扫描的前置操作，为了加速，开启多个nmap进行扫描
// 将ip分割输出到对应目录下的txt文件中，同时会打印对应的nmap命令
func TestMain_GetIPsTxtAndNmapCommand(t *testing.T) {
	svc := service.NewService()
	cmd, err := svc.GetIPsTxtAndNmapCommand("/Users/leyouming/company_program/scan_tool/AssetMapping/nmap/0325_2",
		"ips0325_2",
		100,
		true,
	)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(cmd)
}

// 建立种子和子域名之间的对应关系，之前只将子域名存入，没有建立种子和子域名的对应关系
func TestMain_GenerateSeedSudDomainsAssociation(t *testing.T) {
	svc := service.NewService()
	err := svc.GenerateSeedSudDomainsAssociation()
	if err != nil {
		t.Fatal(err)
	}
}

// 将数据库结果输出到excel文件中
func TestMain_OutputExcel(t *testing.T) {
	svc := service.NewService()
	err := svc.OutputExcel("/Users/leyouming/company_program/script/asset_mapping/file/测绘结果3.xlsx")
	if err != nil {
		t.Fatal(err)
	}
}

// 从excel文件中读取组织和种子以及相关的对应关系，并存入数据库
func TestMain_ReadCompanyAndSeed(t *testing.T) {
	svc := service.NewService()
	err := svc.ReadCompanyAndSeed("/Users/leyouming/company_program/script/asset_mapping/file/备案域名查询2.xlsx")
	if err != nil {
		t.Fatal(err)
	}
}

// 将数据库中的种子生成txt文件，每行是一个domain，主要用于OneForAll之类工具的子域名发现
func TestMain_GenerateDomainTxt(t *testing.T) {
	svc := service.NewService()
	err := svc.GenerateDomainTxt("/Users/leyouming/company_program/script/asset_mapping/file/domains.txt", true)
	if err != nil {
		t.Fatal(err)
	}
}

// 读取OneForAll的结果csv文件，读取其中的子域名和对应IP存入数据库
// 输入文件夹路径，会扫描该路径下的所有csv文件
func TestMain_ReadOneForAllResult(t *testing.T) {
	svc := service.NewService()
	err := svc.ReadOneForAllResult("/Users/leyouming/company_program/scan_tool/OneForAll/results")
	if err != nil {
		t.Fatal(err)
	}
}

// 获取数据库中IP的地理位置，通过https://ipinfo.io/的api
// 需要key
func TestMain_GetIPLocalInfo(t *testing.T) {
	svc := service.NewService()
	err := svc.GetIPLocalInfo()
	if err != nil {
		t.Fatal(err)
	}
}

// 使用爬虫爬取所有子域名网址的title
func TestMain_GetAllDomainTitle(t *testing.T) {
	svc := service.NewService()
	err := svc.GetAllDomainTitle(true)
	if err != nil {
		t.Fatal(err)
	}
}

// 处理没有发现子域名的种子，将他自己作为自己的子域名
func TestMain_DealWithNoSubDomainSeed(t *testing.T) {
	svc := service.NewService()
	err := svc.DealWithNoSubDomainSeed()
	if err != nil {
		t.Fatal(err)
	}
}
