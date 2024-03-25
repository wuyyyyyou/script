package service

import (
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/wuyyyyyou/go-share/pd"

	"github.com/wuyyyyou/script/asset_mapping/models"
	"github.com/wuyyyyou/script/asset_mapping/models/nmap"
	"github.com/wuyyyyou/script/asset_mapping/utils"
)

// ReadSeed 从excel文件中读取种子url，并存入数据库
func (s *Service) ReadSeed(filePath string, sheetName string, rowHeads ...string) error {
	if len(rowHeads) == 0 {
		return fmt.Errorf("rowHeads长度为0")
	}

	df := pd.NewDataFrame(sheetName)
	err := df.ReadExcel(filePath)
	if err != nil {
		return err
	}

	for i := 0; i < df.GetLength(); i++ {
		for _, rowHead := range rowHeads {
			url, _ := df.GetValue(i, rowHead)
			url = utils.GetSeedUrl(url)
			if url == "" {
				logrus.Debugf("%d行%s列的url为空", i, rowHead)
				continue
			}

			seed := &models.Seed{
				SeedName: url,
			}
			err := s.UpsertSeed(seed)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GenerateQuakeQuery 根据数据库中的种子url生成Quake的查询语句
func (s *Service) GenerateQuakeQuery() (string, error) {
	seeds, err := s.GetAllSeeds()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	var num int
	for _, seed := range seeds {
		num++
		sb.WriteString(fmt.Sprintf("domain:\"%s\"", seed.SeedName))
		if num != len(seeds) {
			sb.WriteString(" or ")
		}
	}
	logrus.Debugf("生成了%d条url的quake查询语句", len(seeds))
	return sb.String(), nil
}

// ReadNmapXml 读取nmap的xml文件并存入数据库
func (s *Service) ReadNmapXml(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	nmapRun := new(nmap.NmapRun)
	err = xml.Unmarshal(content, nmapRun)
	if err != nil {
		return err
	}

	for _, host := range nmapRun.Hosts {
		if host.Address.Addr == "" {
			continue
		}
		err := s.UpsertIPAndPorts(&host)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadAllNmapXml 读取目录下的所有nmap的xml文件并存入数据库
func (s *Service) ReadAllNmapXml(dirPath string) error {
	filePaths, err := utils.GetNmapXmlFilesFromDir(dirPath)
	if err != nil {
		return err
	}

	logrus.Debugf("找到nmap结果文件:%#v", filePaths)

	for _, filePath := range filePaths {
		err := s.ReadNmapXml(filePath)
		if err != nil {
			logrus.Errorf("读取nmap结果文件%s失败:%v", filePath, err)
			continue
		}
	}
	return nil
}

// ReadSubDomain 从quake的数据中读取子域名存储到数据库
func (s *Service) ReadSubDomain(filePath string) error {
	df := pd.NewDataFrame("Sheet")
	err := df.ReadExcel(filePath)
	if err != nil {
		return err
	}
	for i := 0; i < df.GetLength(); i++ {
		subDomain, _ := df.GetValue(i, "网站Host（域名）")
		if subDomain == "" {
			logrus.Debugf("%d行的域名为空", i)
			continue
		}
		err := s.UpsertSubDomain(&models.Domain{Domain: subDomain})
		if err != nil {
			return err
		}
	}

	return nil
}

// GetIP 获取数据库中所有子域名的IP
func (s *Service) GetIP() error {
	domains, err := s.GetAllDomains()
	if err != nil {
		return err
	}

	for _, domain := range domains {
		IPs, err := net.LookupIP(domain.Domain)
		if err != nil {
			logrus.Infof("获取域名 %s 的IP失败: %s\n", domain.Domain, err)
			continue
		}
		for _, ip := range IPs {
			ipModel := &models.IP{IP: ip.String()}
			err := s.UpsertIPAndDomain(ipModel, domain)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetIPsTxtAndNmapCommand 将ip分割，生成nmap扫描命令和ip列表文件
func (s *Service) GetIPsTxtAndNmapCommand(dirPath string, fileName string,
	chunkCount int, withoutPortsIP bool) (string, error) {
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		return "", err
	}

	if !dirInfo.IsDir() {
		return "", fmt.Errorf("%s不是目录", dirPath)
	}

	ips, err := s.GetAllIPsWithPorts()
	if err != nil {
		return "", err
	}
	if withoutPortsIP {
		ips = lo.Filter(ips, func(ip *models.IP, _ int) bool {
			return len(ip.Ports) == 0
		})
	}

	var sb strings.Builder
	var num int
	ipChunks := lo.Chunk(ips, chunkCount)
	for _, ipChunk := range ipChunks {
		file, err := os.Create(filepath.Join(dirPath, fmt.Sprintf("%s_%d.txt", fileName, num)))
		if err != nil {
			return "", err
		}

		for _, ip := range ipChunk {
			_, err := file.WriteString(ip.IP + "\n")
			if err != nil {
				return "", err
			}
		}

		sb.WriteString(fmt.Sprintf("nmap -p 1-10000 -sT -sV -iL %s -oX %s -Pn\n",
			fmt.Sprintf("%s_%d.txt", fileName, num),
			fmt.Sprintf("%s_%d.xml", fileName, num),
		))

		num++
		_ = file.Close()
	}

	return sb.String(), nil
}

// GenerateSeedSudDomainsAssociation 建立种子和子域名之间的对应关系
func (s *Service) GenerateSeedSudDomainsAssociation() error {
	seeds, err := s.GetAllSeeds()
	if err != nil {
		return err
	}
	for _, seed := range seeds {
		domains, err := s.FindSubDomains(seed)
		if err != nil {
			return err
		}
		err = s.AddSeedDomain(seed, domains)
		if err != nil {
			return err
		}
	}
	return nil
}

// OutputExcel 数据库数据输出到excel文件
func (s *Service) OutputExcel(filePath string) error {
	df := pd.NewDataFrame("Sheet1")
	seeds, err := s.GetAllSeedsWithAll()
	if err != nil {
		return err
	}

	df.SetHeads([]string{
		"种子",
		"域名",
		"IP",
		"Port",
		"端口协议",
		"协议指纹",
	})

	var num int
	for _, seed := range seeds {
		for _, domain := range seed.Domains {
			for _, ip := range domain.IPs {
				for _, port := range ip.Ports {
					_ = df.SetValue(num, "种子", seed.SeedName)
					_ = df.SetValue(num, "域名", domain.Domain)
					_ = df.SetValue(num, "IP", ip.IP)
					_ = df.SetValue(num, "Port", strconv.Itoa(port.Port))
					_ = df.SetValue(num, "端口协议", port.Protocol)
					_ = df.SetValue(num, "协议指纹", port.Finger)
					num++
				}
			}
		}
	}
	err = df.SaveExcel(filePath)
	if err != nil {
		return err
	}

	return nil
}
