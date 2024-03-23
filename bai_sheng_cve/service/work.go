package service

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/wuyyyyyou/go-share/pd"
)

func (s *Service) readExcel(filepath string) (*pd.DataFrame, error) {
	df := pd.NewDataFrame()
	err := df.ReadExcel(filepath)
	return df, err
}

func (s *Service) syncReadExcel(filepath string) (*pd.SyncDataFrame, error) {
	df := pd.NewSyncDataFrame()
	err := df.ReadExcel(filepath)
	return df, err
}

func (s *Service) GetInfoFromDB(filepath string) error {
	df, err := s.readExcel(filepath)
	if err != nil {
		return err
	}
	for i := 0; i < df.GetLength(); i++ {
		cve, _ := df.GetValue(i, 6)
		cve = strings.ToUpper(cve)
		if strings.HasPrefix(cve, "CVE-") {
			cve = cve[4:]
		}

		cveInfo, err := s.findCveInfoFromDB(cve)
		if err != nil {
			logrus.Errorf("findCveInfoFromDB error: %v", err)
		}

		if cveInfo == nil {
			_ = df.SetValue(i, 7, "")
			_ = df.SetValue(i, 8, "")
		} else {
			_ = df.SetValue(i, 7, cveInfo.Description)
			_ = df.SetValue(i, 8, cveInfo.References)
		}
	}

	newFilepath := strings.Replace(filepath, ".xlsx", "-new.xlsx", 1)
	return df.SaveExcel(newFilepath)
}

func (s *Service) GetInfoFromAPI(filepath string) error {
	df, err := s.syncReadExcel(filepath)
	if err != nil {
		return err
	}

	for i := 0; i < df.GetLength(); i++ {
		cve, _ := df.GetValue(i, 6)
		cve = strings.ToUpper(cve)
		if !strings.HasPrefix(cve, "CVE-") {
			continue
		}

		desc, _ := df.GetValue(i, 7)
		if strings.TrimSpace(desc) != "" {
			continue
		}

		s.swg.Add()
		go func(cve string, index int) {
			defer func() {
				if err := recover(); err != nil {
					logrus.Errorf("GetInfoFromAPI Recovered from panic: %v", err)
				}
			}()
			defer s.swg.Done()

			cveInfo, err := s.findCveInfoFromAPI(cve)
			if err != nil {
				logrus.Errorf("findCveInfoFromAPI error: %v", err)
				return
			}

			if cveInfo == nil {
				_ = df.SetValue(index, 7, "")
				_ = df.SetValue(index, 8, "")
			} else {
				_ = df.SetValue(index, 7, cveInfo.Description)
				_ = df.SetValue(index, 8, cveInfo.References)
			}

		}(cve, i)

	}
	s.swg.Wait()
	newFilepath := strings.Replace(filepath, ".xlsx", "-new.xlsx", 1)
	return df.SaveExcel(newFilepath)
}

func (s *Service) findCveInfoFromAPI(cve string) (*CveInfo, error) {
	var cveInfo *CveInfo
	load, ok := s.CveInfoCacheSync.Load(cve)
	if !ok {
		client := resty.New()

		if s.HttpProxy != "" {
			client.SetProxy(s.HttpProxy)
		}

		responseCveInfo := new(ResponseCveInfo)
		resp, err := client.R().
			SetResult(responseCveInfo).
			Get(fmt.Sprintf("https://v1.cveapi.com/%s.json", cve))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode() != 200 {
			return nil, fmt.Errorf("status code: %d", resp.StatusCode())
		}
		cveInfo = responseCveInfo.ToCevInfo()
		s.CveInfoCacheSync.Store(cve, cveInfo)
	} else {
		cveInfo = load.(*CveInfo)
	}

	return cveInfo, nil
}

func (s *Service) findCveInfoFromDB(cve string) (*CveInfo, error) {
	cveInfo, ok := s.CveInfoCache[cve]
	if !ok {
		var exploitRecords []*ExploitRecord
		err := s.DB.Where("cve like ?", "%"+cve+"%").Find(&exploitRecords).Error
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(exploitRecords); i++ {
			if containsChinese(exploitRecords[i].Description) || i == len(exploitRecords)-1 {
				cveInfo, err = exploitRecords[i].ToCveInfo()
				if err != nil {
					return nil, err
				}
			}
		}

		s.CveInfoCache[cve] = cveInfo
	}

	return cveInfo, nil
}

func (s *Service) TranslateDescription(filepath string) error {
	df, err := s.syncReadExcel(filepath)
	if err != nil {
		return err
	}

	for i := 0; i < df.GetLength(); i++ {
		cve, _ := df.GetValue(i, 6)
		cve = strings.ToUpper(cve)
		desc, _ := df.GetValue(i, 7)
		if !strings.HasPrefix(cve, "CVE-") || strings.TrimSpace(desc) == "" || containsChinese(desc) {
			continue
		}

		s.swg.Add()
		go func(cve string, desc string, index int) {
			defer func() {
				if err := recover(); err != nil {
					logrus.Errorf("GetInfoFromAPI Recovered from panic: %v", err)
				}
			}()
			defer s.swg.Done()

			translateDesc, err := s.translateFromAPI(cve, desc, index)
			if err != nil {
				logrus.Errorf("translateFromAPI error: %v", err)
				return
			}

			_ = df.SetValue(index, 7, translateDesc)
		}(cve, desc, i)
	}

	s.swg.Wait()
	newFilepath := strings.Replace(filepath, ".xlsx", "-new.xlsx", 1)
	return df.SaveExcel(newFilepath)
}

func (s *Service) translateFromAPI(cve string, desc string, index int) (string, error) {
	var translatedDesc string
	load, ok := s.CveDescCacheSync.Load(cve)
	if !ok {
		client := resty.New()

		// 翻译API不需要代理
		//if s.HttpProxy != "" {
		//	client.SetProxy(s.HttpProxy)
		//}

		responseTranslate := new(ResponseTranslate)
		resp, err := client.R().
			SetHeader("Content-Type", "application/x-www-form-urlencoded").
			SetHeader("Authorization", "DeepL-Auth-Key c881c91b-c4e3-449d-82d7-7a3763c47aa4:dp").
			SetFormData(map[string]string{
				"text":        desc,
				"target_lang": "ZH",
			}).
			SetResult(responseTranslate).
			Post("https://api.deepl-pro.com/v2/translate")
		if err != nil {
			return "", err
		}
		if resp.StatusCode() != 200 {
			return "", fmt.Errorf("status code: %d", resp.StatusCode())
		}
		translatedDesc = responseTranslate.Translations[0].Text
		s.CveDescCacheSync.Store(cve, translatedDesc)
	} else {
		translatedDesc = load.(string)
	}

	return translatedDesc, nil
}

func containsChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

func (s *Service) GetCVSSFromAPI(filepath string) error {
	df, err := s.syncReadExcel(filepath)
	if err != nil {
		return err
	}

	for i := 0; i < df.GetLength(); i++ {
		cve, _ := df.GetValue(i, "CVE漏洞编号")
		cve = strings.ToUpper(cve)
		if !strings.HasPrefix(cve, "CVE-") {
			continue
		}

		score, _ := df.GetValue(i, "CVSS评分")
		rank, _ := df.GetValue(i, "CVSS等级")
		if strings.TrimSpace(score) != "" && strings.TrimSpace(rank) != "" {
			continue
		}

		s.swg.Add()
		go func(cve string, index int) {
			defer func() {
				if err := recover(); err != nil {
					logrus.Errorf("GetInfoFromAPI Recovered from panic: %v", err)
				}
			}()
			defer s.swg.Done()

			cvssInfo, err := s.findCveCvssFromAPI(cve)
			if err != nil {
				logrus.Errorf("findCveCvssFromAPI error: %v", err)
				return
			}

			if cvssInfo == nil {
				_ = df.SetValue(index, "CVSS评分", "")
				_ = df.SetValue(index, "CVSS等级", "")
			} else {
				_ = df.SetValue(index, "CVSS评分", cvssInfo.BaseScore)
				_ = df.SetValue(index, "CVSS等级", cvssInfo.BaseSeverity)
			}

		}(cve, i)
	}

	s.swg.Wait()
	newFilepath := strings.Replace(filepath, ".xlsx", "-new.xlsx", 1)
	return df.SaveExcel(newFilepath)
}

func (s *Service) findCveCvssFromAPI(cve string) (*CvssInfo, error) {
	var cvssInfo *CvssInfo
	load, ok := s.CvssInfoCacheSync.Load(cve)
	if !ok {
		client := resty.New().
			SetRetryCount(3).
			SetRetryWaitTime(1 * time.Second).
			SetRetryMaxWaitTime(5 * time.Second)

		if s.HttpProxy != "" {
			client.SetProxy(s.HttpProxy)
		}

		responseCveInfo := new(ResponseCveInfo)
		resp, err := client.R().
			SetResult(responseCveInfo).
			Get(fmt.Sprintf("https://v1.cveapi.com/%s.json", cve))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode() != 200 {
			return nil, fmt.Errorf("status code: %d", resp.StatusCode())
		}
		cvssInfo = responseCveInfo.ToCvssInfo()
		s.CveInfoCacheSync.Store(cve, cvssInfo)
	} else {
		cvssInfo = load.(*CvssInfo)
	}

	return cvssInfo, nil
}
