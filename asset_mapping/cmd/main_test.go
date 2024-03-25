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
		&models.Seed{},
		//&models.IP{},
		//&models.Port{},
		&models.Domain{},
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
	query, err := svc.GenerateQuakeQuery()
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
	err := svc.ReadAllNmapXml("/Users/leyouming/company_program/scan_tool/AssetMapping/nmap/0324")
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
	err := svc.GetIP()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain_GetIPsTxtAndNmapCommand(t *testing.T) {
	svc := service.NewService()
	cmd, err := svc.GetIPsTxtAndNmapCommand("/Users/leyouming/company_program/scan_tool/AssetMapping/nmap/0324",
		"ips0324",
		3,
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
	err := svc.OutputExcel("/Users/leyouming/company_program/script/asset_mapping/file/测绘结果.xlsx")
	if err != nil {
		t.Fatal(err)
	}
}
