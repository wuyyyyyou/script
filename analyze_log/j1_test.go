package analyze_log

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/xuri/excelize/v2"
)

func Test1(t *testing.T) {
	js := `
	{
	"AssetID": "000000cu_ryy02012023062800016",
	"AssetIPInfo": [{
		"IpType": 0,
		"Ip": "10.4.112.130",
		"OpenPort": [{
			"protocol": "",
			"port": 65535
		}]
	}],
	"AssetInfo": [{
		"Brand": [{
			"Manufacturer": "CS0012",
			"Region": "1",
			"Name": "H3C"
		}],
		"Language": "english",
		"Memory": [{
			"Manufacturer": "",
			"Model": ""
		}],
		"OpenSource": "",
		"HardwareEnvironment": "x64",
		"CPU": [{
			"Manufacturer": "CS0342",
			"Model": "XLP308 Rev B2   FPU"
		}],
		"information": [{
			"Version": "Release 9141P42",
			"Model": "XH1480"
		}],
		"SoftwareEnvironment": "linux",
		"Chip": [{
			"Manufacturer": "",
			"Model": ""
		}]
	}],
	"AssetName": "防火墙",
	"AssetTag": "防火墙",
	"AssetType": "0201",
	"DeviceLevel": "2",
	"FoundTypeList": "5,6",
	"FoundTypeTime": "20230628181838",
	"IsAccess": 0,
	"IspCode": "cu_ryy",
	"Location": "610100",
	"NetPosition": 0,
	"NetworkUnit": "10001105",
	"ObjectName": "联通软研院知识中心",
	"OrgCode": "000000",
	"State": 1,
	"SystemName": "知识中心",
	"TaskID": "testTask-1700463672353"
}`
	asset := Asset{}
	err := json.Unmarshal([]byte(js), &asset)
	if err != nil {
		t.Error(err)
	}
	t.Log("ok")
}

func Test2(t *testing.T) {
	service := NewService()
	err := service.LoadAssets("files/content.log")
	if err != nil {
		t.Error(err)
	}
	t.Log("ok")
}

func Test3(t *testing.T) {
	var a StringType = "123"
	t.Log(a)
}

func Test4(t *testing.T) {
	file := excelize.NewFile()
	index, _ := file.NewSheet("Sheet1")

	file.SetActiveSheet(index)

	// 根据内存中的数据填充文件
	data := [][]interface{}{
		{"Name", "Age", "Gender"},
		{"John", 30, "Male"},
		{"Alice", 25, "Female"},
		{"Bob", 22, "Male"},
	}

	for i, row := range data {
		for j, colCell := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+1)
			err := file.SetCellValue("Sheet1", cell, colCell)
			if err != nil {
				t.Error(err)
			}
		}
	}

	if err := file.SaveAs("test.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func Test5(t *testing.T) {
	a := Asset{}
	a.FillNil()
	t.Log("ok")
}
