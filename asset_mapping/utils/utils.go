package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// GetSeedUrl 对种子url做一个预处理
func GetSeedUrl(url string) string {
	url = strings.TrimSpace(url)

	if strings.HasPrefix(url, "http://") {
		url = strings.Replace(url, "http://", "", 1)
	}

	if strings.HasPrefix(url, "https://") {
		url = strings.Replace(url, "https://", "", 1)
	}

	urlList := strings.Split(url, "/")
	url = urlList[0]

	// 端口号要去掉
	urlList = strings.Split(url, ":")
	if len(urlList) == 2 {
		url = urlList[0]
	}
	return url
}

// GetXmlFilesFromDir 获取目录下的所有nmap扫描的xml文件
func GetXmlFilesFromDir(dirPath string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".xml") && !info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			files = append(files, absPath)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func GetCsvFilesFromDir(dirPath string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".csv") && !info.IsDir() {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			files = append(files, absPath)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
