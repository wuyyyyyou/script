package test

import (
	"testing"

	"github.com/go-resty/resty/v2"

	"github.com/wuyyyyou/script/bai_sheng_cve/service"
)

func TestService_NewService(t *testing.T) {
	svc := service.NewService()
	var exploitRecords []service.ExploitRecord
	err := svc.DB.Where("cve like ?", "%"+"2019-0221"+"%").Find(&exploitRecords).Error
	if err != nil {
		t.Error(err)
	}
	t.Log(exploitRecords)
}

func TestService_DeeplAPI(t *testing.T) {
	client := resty.New()
	responseTranslate := new(service.ResponseTranslate)
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Authorization", "DeepL-Auth-Key c881c91b-c4e3-449d-82d7-7a3763c47aa4:dp").
		SetFormData(map[string]string{
			"text":        "In Apache Tomcat 9.0.0.M1 to 9.0.30, 8.5.0 to 8.5.50 and 7.0.0 to 7.0.99 the HTTP header parsing code used an approach to end-of-line parsing that allowed some invalid HTTP headers to be parsed as valid. This led to a possibility of HTTP Request Smuggling if Tomcat was located behind a reverse proxy that incorrectly handled the invalid Transfer-Encoding header in a particular manner. Such a reverse proxy is considered unlikely.",
			"target_lang": "ZH",
		}).
		SetResult(responseTranslate).
		Post("https://api.deepl-pro.com/v2/translate")

	if err != nil {
		t.Error(err)
	}
	t.Log(resp.StatusCode())
	t.Log(&responseTranslate)
}
