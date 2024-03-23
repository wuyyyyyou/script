package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/remeh/sizedwaitgroup"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Service struct {
	CveInfoCache      map[string]*CveInfo
	CveInfoCacheSync  sync.Map
	CveDescCacheSync  sync.Map
	CvssInfoCacheSync sync.Map
	DB                *gorm.DB
	HttpMaxGoroutines int
	HttpProxy         string
	swg               sizedwaitgroup.SizedWaitGroup
}

func NewService() *Service {
	service := new(Service)

	service.CveInfoCache = make(map[string]*CveInfo)
	service.DB = getDB()
	// 设置http访问最大并发数
	service.HttpMaxGoroutines = 10
	// 设置http代理，""则不使用代理
	service.HttpProxy = "http://localhost:7890"

	service.swg = sizedwaitgroup.New(service.HttpMaxGoroutines)
	return service
}

type CveInfo struct {
	Description string
	References  string
}

type CvssInfo struct {
	BaseScore    string
	BaseSeverity string
}

type ExploitRecord struct {
	Cve         string
	Description string
	Url         string
}

type ResponseCveInfo struct {
	Cve    Cve    `json:"cve"`
	Impact Impact `json:"impact"`
}

type Impact struct {
	BaseMetricV3 BaseMetricV3 `json:"baseMetricV3"`
	BaseMetricV2 BaseMetricV2 `json:"baseMetricV2"`
}

type BaseMetricV3 struct {
	CvssV3 *CvssV3 `json:"cvssV3"`
}

type CvssV3 struct {
	BaseScore    float64 `json:"baseScore"`
	BaseSeverity string  `json:"baseSeverity"`
}

type BaseMetricV2 struct {
	CvssV2   *CvssV2 `json:"cvssV2"`
	Severity *string `json:"severity"`
}

type CvssV2 struct {
	BaseScore float64 `json:"baseScore"`
}

type Cve struct {
	References  References  `json:"references"`
	Description Description `json:"description"`
}

type References struct {
	ReferenceData []ReferenceData `json:"reference_data"`
}

type ReferenceData struct {
	Url string `json:"url"`
}

type Description struct {
	DescriptionData []DescriptionData `json:"description_data"`
}

type ResponseTranslate struct {
	Translations []Translation `json:"translations"`
}

type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

type DescriptionData struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

func (r *ResponseCveInfo) ToCevInfo() *CveInfo {
	cveInfo := new(CveInfo)
	var urlList []string
	for _, reference := range r.Cve.References.ReferenceData {
		urlList = append(urlList, reference.Url)
	}
	cveInfo.References = strings.Join(urlList, "\n")

	for _, description := range r.Cve.Description.DescriptionData {
		if description.Lang == "zh" {
			cveInfo.Description = description.Value
		}

		if description.Lang == "en" && cveInfo.Description == "" {
			cveInfo.Description = description.Value
		}
	}

	return cveInfo
}

func (r *ResponseCveInfo) ToCvssInfo() *CvssInfo {
	cvssInfo := new(CvssInfo)
	if r.Impact.BaseMetricV3.CvssV3 != nil {
		cvssInfo.BaseScore = fmt.Sprintf("%.1f", r.Impact.BaseMetricV3.CvssV3.BaseScore)
		cvssInfo.BaseSeverity = r.Impact.BaseMetricV3.CvssV3.BaseSeverity

		return cvssInfo
	}

	if r.Impact.BaseMetricV2.CvssV2 != nil {
		cvssInfo.BaseScore = fmt.Sprintf("%.1f", r.Impact.BaseMetricV2.CvssV2.BaseScore)
		if r.Impact.BaseMetricV2.Severity != nil {
			cvssInfo.BaseSeverity = *r.Impact.BaseMetricV2.Severity
		}

		return cvssInfo
	}

	return cvssInfo
}

func (e *ExploitRecord) ToCveInfo() (*CveInfo, error) {
	cveInfo := new(CveInfo)
	cveInfo.Description = e.Description
	var urlList []string
	err := json.Unmarshal([]byte(e.Url), &urlList)
	if err != nil {
		return nil, err
	}
	cveInfo.References = strings.Join(urlList, "\n")
	return cveInfo, nil
}

func getDB() *gorm.DB {
	username := "root"
	password := "123456"
	host := "127.0.0.1"
	port := 3306
	Dbname := "cybersec_diglab"
	timeout := "10s"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s",
		username, password, host, port, Dbname, timeout)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}

	return db
}
