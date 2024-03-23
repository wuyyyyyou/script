package analyze_log

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

func (s *Service) LoadAssets(jsonPath string) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		return err
	}
	defer CloseAll(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		asset := Asset{}
		err = json.Unmarshal([]byte(line), &asset)
		if err != nil {
			return err
		}
		s.assets = append(s.assets, asset)
	}

	return nil
}

func (s *Service) AssetsToXlsx(xlsxPath string) error {
	data := make([][]any, 0)
	data = append(data, generateHead())
	for count, asset := range s.assets {
		if count%100 == 0 {
			log.Printf("正在处理第 %v 条数据", count)
		}

		asset.FillNil()

		for _, assetIPInfo := range asset.AssetIPInfos {
			for _, openPort := range assetIPInfo.OpenPorts {
				for _, assetInfo := range asset.AssetInfos {
					for _, brand := range assetInfo.Brands {
						for _, memory := range assetInfo.Memorys {
							for _, cpu := range assetInfo.CPUs {
								for _, information := range assetInfo.Informations {
									for _, chip := range assetInfo.Chips {
										row := []any{
											asset.AssetID,
											assetIPInfo.IPType,
											assetIPInfo.IP,
											openPort.Protocol,
											openPort.Port,
											brand.Manufacturer,
											brand.Region,
											brand.Name,
											assetInfo.Language,
											memory.Manufacturer,
											memory.Model,
											assetInfo.OpenSource,
											assetInfo.HardwareEnvironment,
											cpu.Manufacturer,
											cpu.Model,
											information.Version,
											information.Model,
											assetInfo.SoftwareEnvironment,
											chip.Manufacturer,
											chip.Model,
											asset.AssetName,
											asset.AssetTag,
											asset.AssetType,
											asset.DeviceLevel,
											asset.FoundTypeList,
											asset.FoundTypeTime,
											asset.IsAccess,
											asset.IspCode,
											asset.Location,
											asset.NetPosition,
											asset.NetworkUnit,
											asset.ObjectName,
											asset.OrgCode,
											asset.State,
											asset.SystemName,
											asset.TaskID,
										}
										data = append(data, row)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return saveDataToXlsx(xlsxPath, data)
}

func saveDataToXlsx(xlsxPath string, data [][]any) error {
	file := excelize.NewFile()
	index, _ := file.NewSheet("Sheet1")
	file.SetActiveSheet(index)

	for i, row := range data {
		if i%100 == 0 {
			log.Printf("正在写入第 %v 条数据", i)
		}
		for j, colCell := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+1)
			err := file.SetCellValue("Sheet1", cell, colCell)
			if err != nil {
				return err
			}
		}
	}

	if err := file.SaveAs(xlsxPath); err != nil {
		return err
	}

	return nil
}

func generateHead() []any {
	return []any{
		"AssetID",
		"IpType",
		"Ip",
		"protocol",
		"port",
		"BrandManufacturer",
		"BrandRegion",
		"BrandName",
		"Language",
		"MemoryManufacturer",
		"MemoryModel",
		"OpenSource",
		"HardwareEnvironment",
		"CPUManufacturer",
		"CPUModel",
		"informationVersion",
		"informationModel",
		"SoftwareEnvironment",
		"ChipManufacturer",
		"ChipModel",
		"AssetName",
		"AssetTag",
		"AssetType",
		"DeviceLevel",
		"FoundTypeList",
		"FoundTypeTime",
		"IsAccess",
		"IspCode",
		"Location",
		"NetPosition",
		"NetworkUnit",
		"ObjectName",
		"OrgCode",
		"State",
		"SystemName",
		"TaskID",
	}
}
