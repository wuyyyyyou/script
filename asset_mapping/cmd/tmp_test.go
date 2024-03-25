package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/wuyyyyou/script/asset_mapping/models"
	"github.com/wuyyyyou/script/asset_mapping/service"
)

func TestTmp_Updates(t *testing.T) {
	svc := service.NewService()
	seed := &models.Seed{
		SeedName: "asi-callback.shizhuang-inc.com",
	}
	err := svc.UpsertSeed(seed)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", seed)

}

func TestTmp_NoInfoIP(t *testing.T) {
	svc := service.NewService()
	bytes, err := os.ReadFile("/Users/leyouming/company_program/scan_tool/AssetMapping/nmap/ips20240322.txt")
	if err != nil {
		t.Fatal(err)
	}

	ips := strings.Split(string(bytes), "\n")
	logrus.Debugf("共有%d个ip", len(ips))
	for _, ip := range ips {
		var infoIps []*models.IP
		svc.DB.Where(&models.IP{IP: ip}).First(&infoIps)
		if len(infoIps) == 0 {
			logrus.Debugf("ip:%s 无信息", ip)
		}
	}
}

func TestTmp_1(t *testing.T) {
	s := "a.csv.bb.csv"
	s = strings.Replace(s, ".csv", "", 1)
	t.Log(s)
}
