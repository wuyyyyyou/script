package utils

import (
	"fmt"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/wuyyyyyou/go-share/pd"
)

// GetQuakeQuery 读取excel文件，获取域名，生成Quake的查询语句
func GetQuakeQuery(filepath string) string {
	df := pd.NewDataFrame("域名URL")
	//err := df.ReadExcel("/Users/leyouming/program/golang_note/golang_code/go_code/script/company_domain/file/测绘2.xlsx")
	err := df.ReadExcel(filepath)
	if err != nil {
		logrus.Fatal(err)
	}

	cache := make(map[string]string)
	for i := 0; i < df.GetLength(); i++ {
		url1, _ := df.GetValue(i, "域名/URL-1")
		url2, _ := df.GetValue(i, "域名/URL-2")

		url1 = preprocessUrl(url1)
		url2 = preprocessUrl(url2)

		if url1 != "" {
			setUrlQueryCache(url1, cache)
		}
		if url2 != "" {
			setUrlQueryCache(url2, cache)
		}

	}

	var sb strings.Builder
	var num int
	for _, urlQuery := range cache {
		num++
		sb.WriteString(urlQuery)
		if num != len(cache) {
			sb.WriteString(" or ")
		}
	}

	return sb.String()
}

func setUrlQueryCache(url string, cache map[string]string) {
	_, ok := cache[url]
	if !ok {
		cache[url] = fmt.Sprintf("domain:\"%s\"", url)
	} else {
		fmt.Printf("url: %s 已经存在\n", url)
	}
}

func preprocessUrl(url string) string {
	url = strings.TrimSpace(url)

	if strings.HasPrefix(url, "http://") {
		url = strings.Replace(url, "http://", "", 1)
	}

	if strings.HasPrefix(url, "https://") {
		url = strings.Replace(url, "https://", "", 1)
	}

	urlList := strings.Split(url, "/")
	url = urlList[0]

	return url
}

func GetIP(input string, output string) error {
	df := pd.NewDataFrame("Sheet1")
	err := df.ReadExcel(input)
	if err != nil {
		return err
	}

	dfOut := pd.NewDataFrame("Sheet1")
	dfOut.SetHeads([]string{
		"域名",
		"IP",
	})

	indexOut := 0
	for i := 0; i < df.GetLength(); i++ {
		domain, _ := df.GetValue(i, "网站Host（域名）")
		IPs, err := net.LookupIP(domain)
		if err != nil {
			logrus.Infof("获取域名 %s 的IP失败: %s\n", domain, err)
			continue
		}
		for _, ip := range IPs {
			_ = dfOut.SetValue(indexOut, "域名", domain)
			_ = dfOut.SetValue(indexOut, "IP", ip.String())
			indexOut++
		}
	}

	err = dfOut.SaveExcel(output)
	if err != nil {
		return err
	}
	return nil
}
