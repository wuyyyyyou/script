package cmd

import (
	"fmt"
	"testing"

	"github.com/wuyyyyou/script/asset_mapping/models"
	"github.com/wuyyyyou/script/asset_mapping/service"
)

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

func TestMain_GenerateQuakeQuery(t *testing.T) {
	svc := service.NewService()
	query, err := svc.GenerateQuakeQuery(true)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("'%s'", query)
}

func TestMain_ReadNmapXml(t *testing.T) {
	svc := service.NewService()
	err := svc.ReadNmapXml("/Users/leyouming/company_program/script/asset_mapping/file/test1.xml")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain_ReadAllNmapXml(t *testing.T) {
	svc := service.NewService()
	err := svc.ReadAllNmapXml("/Users/leyouming/company_program/scan_tool/AssetMapping/nmap/0325_1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain_ReadSubDomain(t *testing.T) {
	svc := service.NewService()
	err := svc.ReadSubDomain("/Users/leyouming/company_program/script/asset_mapping/file/服务数据_20240324_191828.xlsx")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain_GetIP(t *testing.T) {
	svc := service.NewService()
	err := svc.GetIP(true)
	if err != nil {
		t.Fatal(err)
	}
}

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

func TestMain_GenerateSeedSudDomainsAssociation(t *testing.T) {
	svc := service.NewService()
	err := svc.GenerateSeedSudDomainsAssociation()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain_OutputExcel(t *testing.T) {
	svc := service.NewService()
	err := svc.OutputExcel("/Users/leyouming/company_program/script/asset_mapping/file/测绘结果2.xlsx")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain_ReadCompanyAndSeed(t *testing.T) {
	svc := service.NewService()
	err := svc.ReadCompanyAndSeed("/Users/leyouming/company_program/script/asset_mapping/file/备案域名查询2.xlsx")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain_GenerateDomainTxt(t *testing.T) {
	svc := service.NewService()
	err := svc.GenerateDomainTxt("/Users/leyouming/company_program/script/asset_mapping/file/domains.txt", true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain_ReadOneForAllResult(t *testing.T) {
	svc := service.NewService()
	err := svc.ReadOneForAllResult("/Users/leyouming/company_program/scan_tool/OneForAll/results")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain_GetIPLocalInfo(t *testing.T) {
	svc := service.NewService()
	err := svc.GetIPLocalInfo()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain_GetAllDomainTitle(t *testing.T) {
	svc := service.NewService()
	err := svc.GetAllDomainTitle(true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain_DealWithNoSubDomainSeed(t *testing.T) {
	svc := service.NewService()
	err := svc.DealWithNoSubDomainSeed()
	if err != nil {
		t.Fatal(err)
	}
}
