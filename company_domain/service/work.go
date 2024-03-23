package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/wuyyyyyou/go-share/pd"
	"gorm.io/gorm"
)

func (s *Service) WriteDomainDao(filePath string, sheetName string) error {
	df := pd.NewDataFrame(sheetName)
	err := df.ReadExcel(filePath)
	if err != nil {
		return err
	}

	for i := 0; i < df.GetLength(); i++ {
		companyName, _ := df.GetValue(i, "企业名称")
		comment, _ := df.GetValue(i, "备注")
		companyType, _ := df.GetValue(i, "类别")

		companyName = strings.TrimSpace(companyName)
		comment = strings.TrimSpace(comment)
		companyType = strings.TrimSpace(companyType)

		if companyName == "" {
			continue
		}

		err := s.DB.Transaction(func(tx *gorm.DB) error {
			var count int64
			err := tx.Model(&CompanyInfo{}).Where("company_name = ?", companyName).Count(&count).Error
			if err != nil {
				return err
			}

			if count > 0 {
				logrus.Infof("公司:'%s'已经存在", companyName)
				return nil
			}

			companyInfo := &CompanyInfo{
				CompanyName: companyName,
				Comment:     comment,
				Type:        companyType,
			}
			err = tx.Create(companyInfo).Error
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetICPDomainInfo() error {
	companyInfos := make([]*CompanyInfo, 0)
	err := s.DB.Where("finished = ?", 0).Find(&companyInfos).Error
	if err != nil {
		return err
	}
	logrus.Infof("读取了%d条未完成查询的企业信息", len(companyInfos))

	for _, companyInfo := range companyInfos {
		s.swg.Add()
		go func(companyInfo *CompanyInfo) {
			defer func() {
				if r := recover(); r != nil {
					logrus.Errorf("Recovered from panic: %v\n", r)
				}
			}()
			defer s.swg.Done()

			err := s.GetICPDomainInfoFromAPI(companyInfo)
			if err != nil {
				logrus.Errorf("获取公司'%s'的备案信息失败:'%v'", companyInfo.CompanyName, err)
			}

		}(companyInfo)
	}

	return nil

}

// GetICPDomainInfoFromAPI 通过站长之家获取网站备案信息并存储
func (s *Service) GetICPDomainInfoFromAPI(companyInfo *CompanyInfo) error {
	client := resty.New()
	domainResponse := new(DomainResponse)
	resp, err := client.R().
		SetResult(domainResponse).
		SetQueryParam("key", s.ChinazKey).
		SetQueryParam("companyname", companyInfo.CompanyName).
		Get("https://apidatav2.chinaz.com/single/sponsorunit")

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("请求状态码:%d", resp.StatusCode())
	}

	if domainResponse.StateCode != 1 {
		return fmt.Errorf("站长之家状态码:%d", domainResponse.StateCode)
	}

	for _, result := range domainResponse.Results {
		result.BelongToCompany = companyInfo.ID
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&CompanyInfo{}).Where("id = ?", companyInfo.ID).Update("finished", 1).Error
		if err != nil {
			return err
		}

		err = tx.Create(&domainResponse.Results).Error
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

func (s *Service) ReadJsonAndWriteDomain(companyName string) error {
	j, err := os.ReadFile("/Users/leyouming/program/golang_note/golang_code/go_code/script/company_domain/file/tmp.json")
	if err != nil {
		return err
	}

	domainResponse := new(DomainResponse)
	err = json.Unmarshal(j, domainResponse)
	if err != nil {
		return err
	}

	err = s.updateDomainInfo(companyName, domainResponse)
	if err != nil {
		return err
	}

	return nil

}

func (s *Service) updateDomainInfo(companyName string, domainResponse *DomainResponse) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		companyInfo := new(CompanyInfo)
		err := tx.Where("company_name = ?", companyName).First(companyInfo).Error
		if err != nil {
			return err
		}

		err = tx.Model(&CompanyInfo{}).Where("id = ?", companyInfo.ID).Update("finished", 1).Error
		if err != nil {
			return err
		}

		for _, result := range domainResponse.Results {
			result.BelongToCompany = companyInfo.ID
		}
		err = tx.Create(&domainResponse.Results).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) ReadTxtAndWriteDomain(companyName string) error {
	content, err := os.ReadFile("/Users/leyouming/program/golang_note/golang_code/go_code/script/company_domain/file/tmp.txt")
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines)%6 != 0 {
		return fmt.Errorf("文件行数不是6的倍数")
	}

	domainResponse := new(DomainResponse)
	results := make([]*DomainICPInfo, 0)
	for i := 0; i < len(lines); i += 6 {
		domain := &DomainICPInfo{
			ID:              0,
			BelongToCompany: 0,
			Domain:          lines[i+4],
			Owner:           "",
			CompanyType:     "",
			SiteLicense:     lines[i+1],
			SiteName:        lines[i+3],
			MainPage:        "www." + lines[i+4],
			VerifyTime:      lines[i+5],
		}
		results = append(results, domain)
	}
	domainResponse.Results = results

	err = s.updateDomainInfo(companyName, domainResponse)
	if err != nil {
		return err
	}

	return nil
}

// WriteExcel 将公司的域名备案信息写入Excel
func (s *Service) WriteExcel(filePath string) error {
	companyInfos := make([]*CompanyInfo, 0)

	err := s.DB.Where("finished = ?", 1).
		Preload("DomainInfos").Find(&companyInfos).Error
	if err != nil {
		return err
	}

	df := pd.NewDataFrame("Sheet1")
	df.SetHeads([]string{
		"序号",
		"企业名称",
		"备注",
		"类别",
		"域名",
		"统一社会信用代码",
	})

	var index int
	var rowIndex int
	for _, companyInfo := range companyInfos {
		index++
		if len(companyInfo.DomainInfos) == 0 {
			_ = df.SetValue(rowIndex, "序号", strconv.Itoa(index))
			_ = df.SetValue(rowIndex, "企业名称", companyInfo.CompanyName)
			_ = df.SetValue(rowIndex, "备注", companyInfo.Comment)
			_ = df.SetValue(rowIndex, "类别", companyInfo.Type)
			_ = df.SetValue(rowIndex, "域名", "")
			_ = df.SetValue(rowIndex, "统一社会信用代码", companyInfo.Credit)
			rowIndex++
		} else {
			for _, domain := range companyInfo.DomainInfos {
				_ = df.SetValue(rowIndex, "序号", strconv.Itoa(index))
				_ = df.SetValue(rowIndex, "企业名称", companyInfo.CompanyName)
				_ = df.SetValue(rowIndex, "备注", companyInfo.Comment)
				_ = df.SetValue(rowIndex, "类别", companyInfo.Type)
				_ = df.SetValue(rowIndex, "域名", domain.Domain)
				_ = df.SetValue(rowIndex, "统一社会信用代码", companyInfo.Credit)
				rowIndex++
			}
		}
	}

	err = df.SaveExcel(filePath)
	if err != nil {
		return err
	}
	return nil
}
