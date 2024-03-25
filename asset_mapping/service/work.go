package service

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/wuyyyyyou/go-share/pd"
	"github.com/wuyyyyyou/go-share/share"

	"github.com/wuyyyyou/script/asset_mapping/models"
	"github.com/wuyyyyou/script/asset_mapping/models/dto"
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

// ReadCompanyAndSeed 从excel文件中读取组织和种子以及相关的对应关系，并存入数据库
func (s *Service) ReadCompanyAndSeed(filePath string) error {
	df := pd.NewDataFrame("Sheet1")
	err := df.ReadExcel(filePath)
	if err != nil {
		return err
	}

	for i := 0; i < df.GetLength(); i++ {
		company, _ := df.GetValue(i, "企业名称")
		seed, _ := df.GetValue(i, "域名")
		credit, _ := df.GetValue(i, "统一社会信用代码")

		company = strings.TrimSpace(company)
		seed = strings.TrimSpace(seed)

		if company == "" || seed == "" {
			logrus.Debugf("第%d行公司'%s'或者种子'%s'为空", i, company, seed)
			continue
		}

		companyModel := &models.Company{
			Company: company,
			Credit:  credit,
		}
		seedModel := &models.Seed{
			SeedName: seed,
		}

		err := s.UpsertCompanyAndSeed(companyModel, seedModel)
		if err != nil {
			return err
		}
	}

	return nil
}

// GenerateQuakeQuery 根据数据库中的种子url生成Quake的查询语句
func (s *Service) GenerateQuakeQuery(withoutSubdomain bool) (string, error) {
	var seeds []*models.Seed
	var err error
	if withoutSubdomain {
		seeds, err = s.GetAllSeedsWithDomains()
		if err != nil {
			return "", err
		}
		seeds = lo.Filter(seeds, func(seed *models.Seed, _ int) bool {
			return len(seed.Domains) == 0
		})
	} else {
		seeds, err = s.GetAllSeeds()
		if err != nil {
			return "", err
		}
	}

	var sb strings.Builder
	var num int
	for _, seed := range seeds {
		num++
		sb.WriteString(fmt.Sprintf("domain:\"%s\"", seed.SeedName))
		if num != len(seeds) {
			sb.WriteString(" || ")
		}
	}
	logrus.Debugf("生成了%d条url的quake查询语句", len(seeds))
	return sb.String(), nil
}

// GenerateDomainTxt 将数据库中的domain生成txt文件，每行是一个domain
func (s *Service) GenerateDomainTxt(filePath string, withoutSubdomain bool) error {
	var seeds []*models.Seed
	var err error
	if withoutSubdomain {
		seeds, err = s.GetAllSeedsWithDomains()
		if err != nil {
			return err
		}
		seeds = lo.Filter(seeds, func(seed *models.Seed, _ int) bool {
			return len(seed.Domains) == 0
		})
	} else {
		seeds, err = s.GetAllSeeds()
		if err != nil {
			return err
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	for _, seed := range seeds {
		_, err := file.WriteString(seed.SeedName + "\n")
		if err != nil {
			return err
		}
	}
	return nil
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
	filePaths, err := utils.GetXmlFilesFromDir(dirPath)
	if err != nil {
		return err
	}

	logrus.Debugf("找到%d个nmap结果文件:%#v", len(filePaths), filePaths)

	for _, filePath := range filePaths {
		err := s.ReadNmapXml(filePath)
		if err != nil {
			logrus.Errorf("读取nmap结果文件%s失败:%v", filePath, err)
			continue
		}
	}
	return nil
}

// ReadSubDomain 从quake的excel数据中读取子域名存储到数据库
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
		err := s.UpsertDomain(&models.Domain{Domain: subDomain})
		if err != nil {
			return err
		}
	}

	return nil
}

// GetIP 获取数据库中所有子域名的IP
func (s *Service) GetIP(withoutIP bool) error {
	domains, err := s.GetAllDomainsWithIPs()
	if err != nil {
		return err
	}

	if withoutIP {
		domains = lo.Filter(domains, func(domain *models.Domain, _ int) bool {
			return len(domain.IPs) == 0
		})
	}

	for _, domain := range domains {
		IPs, err := net.LookupIP(domain.Domain)
		if err != nil {
			logrus.Infof("获取域名 %s 的IP失败: %s\n", domain.Domain, err)
			continue
		}
		for _, ip := range IPs {
			ipModel := &models.IP{IP: ip.String()}
			err := s.UpsertIPWithDomain(ipModel, domain)
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

		sb.WriteString(fmt.Sprintf("nmap -F -sT -sV -iL %s -oX %s -Pn\n",
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

// ReadOneForAllResult 读取OneForAll的结果csv文件，读取其中的子域名和对应IP存入数据库
func (s *Service) ReadOneForAllResult(dirPath string) error {
	files, err := utils.GetCsvFilesFromDir(dirPath)
	if err != nil {
		return err
	}

	logrus.Debugf("找到OneForAll结果文件:%#v", files)

	for _, filePath := range files {
		df := pd.NewDataFrame()
		err := df.ReadCsv(filePath)
		if err != nil {
			return err
		}

		for i := 0; i < df.GetLength(); i++ {
			subDomain, _ := df.GetValue(i, "subdomain")
			ipStr, _ := df.GetValue(i, "ip")
			title, _ := df.GetValue(i, "title")

			subDomain = strings.TrimSpace(subDomain)
			ipStr = strings.TrimSpace(ipStr)

			if subDomain == "" || ipStr == "" {
				logrus.Debugf("第%d行的子域名'%s'或ip'%s'为空", i, subDomain, ipStr)
			}

			domainModel := &models.Domain{
				Domain: subDomain,
				Title:  title,
			}

			ips := strings.Split(ipStr, ",")
			ips = lo.Filter(ips, func(ip string, _ int) bool {
				return len(ip) != 0
			})
			ipsModel := lo.Map(ips, func(ip string, _ int) *models.IP {
				return &models.IP{
					IP: ip,
				}
			})

			err := s.UpsertDomainAndIPs(domainModel, ipsModel)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DealWithNoSubDomainSeed 处理没有发现子域名的种子，将他自己作为自己的子域名
func (s *Service) DealWithNoSubDomainSeed() error {
	seeds, err := s.GetAllSeedsWithDomains()
	if err != nil {
		return err
	}

	seeds = lo.Filter(seeds, func(seed *models.Seed, _ int) bool {
		return len(seed.Domains) == 0
	})

	for _, seed := range seeds {
		domain := &models.Domain{
			Domain: seed.SeedName,
		}
		err := s.UpsertDomain(domain)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetIPLocalInfo 给数据库中的IP查询地址信息
func (s *Service) GetIPLocalInfo() error {
	ips, err := s.GetAllIPsWithoutLocation()
	if err != nil {
		return err
	}

	for _, ip := range ips {
		s.swg.Add()
		go func(ip *models.IP) {
			defer share.Except()
			defer s.swg.Done()

			client := resty.New()
			ipInfo := new(dto.IpInfo)
			resp, err := client.R().
				SetPathParam("ip", ip.IP).
				SetQueryParam("token", s.IPInfoKey).
				SetResult(ipInfo).
				Get("https://ipinfo.io/{ip}")
			if err != nil {
				logrus.Errorf("请求IP信息错误:%s", err)
				return
			}

			if resp.StatusCode() != 200 {
				logrus.Errorf("请求IP信息状态码错误，状态码:%d", resp.StatusCode())
				return
			}

			if ip.IP != ipInfo.IP {
				logrus.Errorf("查询返回IP和查询IP不同")
				return
			}

			newIP := &models.IP{
				IP:       ip.IP,
				Province: ipInfo.Region,
				City:     ipInfo.City,
			}

			err = s.UpsertIP(newIP)
			if err != nil {
				logrus.Errorf("数据库更新IP失败")
				return
			}
		}(ip)
	}
	s.swg.Wait()

	return nil
}

// GetAllDomainTitle 通过爬虫获取网页title
func (s *Service) GetAllDomainTitle(withoutTitle bool) error {
	domains, err := s.GetAllDomains()
	if err != nil {
		return err
	}

	if withoutTitle {
		domains = lo.Filter(domains, func(domain *models.Domain, _ int) bool {
			return domain.Title == ""
		})
	}

	for _, domain := range domains {
		s.swg.Add()
		go func(domain *models.Domain) {
			defer share.Except()
			defer s.swg.Done()

			var title string
			var err error

			load, ok := s.Cache.Load(domain.Domain)
			if ok {
				title = load.(string)
			} else {
				d := "http://" + domain.Domain
				title, err = s.GetDomainTitle(d)
				if err != nil {
					logrus.Errorf("爬取%s出错:%s", d, err)

					d = "https://" + domain.Domain
					title, err = s.GetDomainTitle(d)
					if err != nil {
						logrus.Errorf("爬取%s出错:%s", d, err)
						return
					}
				}

				s.Cache.Store(domain.Domain, title)
			}

			domain.Title = title
			err = s.UpsertDomain(domain)
			if err != nil {
				logrus.Errorf("存储域名%s出错%s", domain.Domain, err)
			}
		}(domain)
	}
	s.swg.Wait()
	return nil
}

func (s *Service) GetDomainTitle(domain string) (string, error) {
	client := resty.New().
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		SetTimeout(5 * time.Second)

	resp, err := client.R().Get(domain)
	if err != nil {
		return "", err
	}

	html, err := share.ConvertEncoding(resp.Header().Get("Content-Type"), resp.Body())
	if err != nil {
		return "", err
	}

	doc := soup.HTMLParse(string(html))
	element := doc.Find("title")
	if element.Error != nil {
		return "", element.Error
	}

	return element.Text(), nil
}

// OutputExcel 数据库数据输出到excel文件
func (s *Service) OutputExcel(filePath string) error {
	df := pd.NewDataFrame("Sheet1")
	seeds, err := s.GetAllSeedsWithAll()
	if err != nil {
		return err
	}

	df.SetHeads([]string{
		"组织",
		"业务系统",
		"种子",
		"域名",
		"IP",
		"端口",
		"端口协议",
		"指纹协议",
		"指纹组件",
		"省",
		"市",
		"操作系统",
		"系统名称",
	})

	var num int
	for _, seed := range seeds {
		// 就算组织没有组织也要输出
		if len(seed.Companys) == 0 {
			seed.Companys = append(seed.Companys, &models.Company{})
		}
		for _, company := range seed.Companys {
			for _, domain := range seed.Domains {
				for _, ip := range domain.IPs {
					for _, port := range ip.Ports {
						_ = df.SetValue(num, "组织", company.Company)
						_ = df.SetValue(num, "种子", seed.SeedName)
						_ = df.SetValue(num, "域名", domain.Domain)
						_ = df.SetValue(num, "IP", ip.IP)
						_ = df.SetValue(num, "端口", strconv.Itoa(port.Port))
						_ = df.SetValue(num, "端口协议", port.Protocol)
						_ = df.SetValue(num, "指纹协议", port.FingerProtocol)
						_ = df.SetValue(num, "指纹组件", port.FingerElement)
						_ = df.SetValue(num, "省", ip.Province)
						_ = df.SetValue(num, "市", ip.City)
						_ = df.SetValue(num, "系统名称", domain.Title)
						num++
					}
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
