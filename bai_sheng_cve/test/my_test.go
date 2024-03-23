package test

import (
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/wuyyyyyou/go-share/pd"

	"github.com/wuyyyyou/script/bai_sheng_cve/service"
)

func TestMy_1(t *testing.T) {
	df := pd.NewDataFrame()
	newdf := pd.NewDataFrame("更新")
	err := df.ReadExcel("../file/百胜-3.13.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
	newdf.SetHeads(df.GetHeads())
	rows := df.GetRows()
	var newRows [][]string
	for _, row := range rows {
		cveStr := row[6]
		cves := strings.Split(cveStr, "\n")
		for _, cve := range cves {
			newRow := make([]string, 6)
			copy(newRow, row[:6])
			newRow = append(newRow, cve)
			newRows = append(newRows, newRow)
		}
	}
	newdf.SetRows(newRows)
	err = newdf.SaveExcel("../file/百胜-3.13-new.xlsx")
}

func TestMy_2(t *testing.T) {
	client := resty.New()

	client.SetProxy("http://localhost:7890")

	responseCveInfo := new(service.ResponseCveInfo)
	resp, err := client.R().
		SetResult(responseCveInfo).
		Get("https://v1.cveapi.com/CVE-2020-1935.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.StatusCode())
	t.Logf("%#v", responseCveInfo)
}

func TestMy_3(t *testing.T) {
	client := resty.New()

	client.SetProxy("http://localhost:7890")

	responseCveInfo := new(service.ResponseCveInfo)
	resp, err := client.R().
		SetResult(responseCveInfo).
		Get("https://v1.cveapi.com/CVE-2020-19351.json")
	if err != nil {
		t.Error(err)
	}
	t.Log(resp.StatusCode())
	t.Logf("%#v", responseCveInfo.ToCvssInfo())
}
