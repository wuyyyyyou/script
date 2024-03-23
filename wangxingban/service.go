package wangxingban

import (
	"context"
	"errors"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/wuyyyyyou/go-share/httputils"
	"github.com/wuyyyyyou/go-share/pd"
	"github.com/wuyyyyyou/go-share/share"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	SrcExcel string
	DstExcel string
	DB       *mongo.Client
	sem      chan struct{}
	wg       sync.WaitGroup

	ipInfoCache     map[string]any
	domainInfoCache sync.Map

	df *pd.SyncDataFrame
}

func NewService(srcExcel string, dstExcel string) *Service {
	client, err := mongo.Connect(context.Background(), options.Client().
		ApplyURI("mongodb://root:root@localhost:27022"))

	if err != nil {
		panic(err)
	}

	return &Service{
		SrcExcel: srcExcel,
		DstExcel: dstExcel,
		DB:       client,
		sem:      make(chan struct{}, 20),

		ipInfoCache: make(map[string]any),
		df:          pd.NewSyncDataFrame(),
	}
}

func (s *Service) Close() {
	err := s.DB.Disconnect(context.Background())
	if err != nil {
		Logger.Error(err)
	}
}

func (s *Service) Worker() error {
	err := s.df.ReadExcel(s.SrcExcel)
	if err != nil {
		return err
	}

	Logger.Debugf("读取文件 %s 完成, 共 %d 行\n", s.SrcExcel, s.df.GetLength())

	for i := 0; i < s.df.GetLength(); i++ {
		ip, err := s.df.GetValue(i, "IP")
		if err != nil {
			Logger.Errorf("获取第 %d 行的 IP:%s 失败: %s\n", i, ip, err)
			continue
		}
		strings.Trim(ip, " ")

		domain, err := s.df.GetValue(i, "域名")
		if err != nil {
			Logger.Errorf("获取第 %d 行的 域名:%s 失败: %s\n", i, domain, err)
			continue
		}
		strings.Trim(domain, " ")

		if i%100 == 0 {
			Logger.Debugf("正在处理第 %d 行, ip: %s, domain: %s\n", i, ip, domain)
		}

		// 去mongo数据库查询IP数据
		ipInfo, ok := s.ipInfoCache[ip]
		if !ok {
			ipInfo, err = s.getIpInfo(ip)
			if err != nil {
				Logger.Errorf("数据库查询第 %d 行的 IP:%s 信息失败: %s\n", i, ip, err)
				continue
			}
			s.ipInfoCache[ip] = ipInfo
		}

		ipMap, _ := ipInfo.(map[string]any)
		ipMap2, _ := ipMap["ip"]
		ipInfoMap, ok := ipMap2.(map[string]any)
		if !ok {
			Logger.Debugf("数据库查询第 %d 行的 IP:%s 信息为空\n", i, ip)
			s.df.SetValue(i, "省", "")
			s.df.SetValue(i, "市", "")
			s.df.SetValue(i, "经度", "")
			s.df.SetValue(i, "纬度", "")
			s.df.SetValue(i, "操作系统", "")
		} else {
			province, ok := ipInfoMap["province"].(string)
			if !ok || province == "" {
				province = ""
			}
			city, ok := ipInfoMap["city"].(string)
			if !ok || city == "" {
				city = ""
			}
			longitude, ok := ipInfoMap["longitude"].(float64)
			var longitudeStr string
			if !ok || longitude == 0 {
				longitudeStr = ""
			} else {
				longitudeStr = strconv.FormatFloat(longitude, 'f', 4, 64)
			}
			latitude, ok := ipInfoMap["latitude"].(float64)
			var latitudeStr string
			if !ok || latitude == 0 {
				latitudeStr = ""
			} else {
				latitudeStr = strconv.FormatFloat(latitude, 'f', 4, 64)
			}
			osType, ok := ipInfoMap["os_type"].(string)
			if !ok || osType == "" {
				osType = ""
			}

			s.df.SetValue(i, "省", province)
			s.df.SetValue(i, "市", city)
			s.df.SetValue(i, "经度", longitudeStr)
			s.df.SetValue(i, "纬度", latitudeStr)
			s.df.SetValue(i, "操作系统", osType)
		}

		// 爬取domain的title
		s.sem <- struct{}{}
		s.wg.Add(1)
		go func(rowIndex int, domain string) {
			defer share.Except()
			defer func() {
				<-s.sem
				s.wg.Done()
			}()
			title := s.getDomainTitle(domain)
			s.df.SetValue(rowIndex, "系统名称", title)
		}(i, domain)

	}
	s.wg.Wait()

	err = s.df.SaveExcel(s.DstExcel)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) getIpInfo(ip string) (map[string]any, error) {
	collection := s.DB.Database("237_asm").Collection("ip_record")
	filter := bson.M{"ip_address": ip}
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.M{
		"ip.province":  1,
		"ip.city":      1,
		"ip.os_type":   1,
		"ip.latitude":  1,
		"ip.longitude": 1,
		"_id":          0,
	})

	var result map[string]any
	err := collection.FindOne(context.Background(), filter, findOptions).Decode(&result)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	return result, nil
}

func (s *Service) getDomainTitle(domain string) string {
	title, ok := s.domainInfoCache.Load(domain)
	if !ok {
		httpRequest := httputils.NewHttpRequest(domain)
		defer httpRequest.Close()
		httpRequest.SkipVerify()
		httpRequest.SetTimeout(3 * time.Second)

		if err := httpRequest.Get(); err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				//title = "timeout"
				title = ""
			} else {
				//title = err.Error()
				title = ""
			}
		} else if httpRequest.GetResponseStatusCode() != 200 {
			//title = fmt.Sprintf("status code: %d", httpRequest.GetResponseStatusCode())
			title = ""
		} else if html, err := httpRequest.GetBodyStringEncoding(); err != nil {
			//title = err.Error()
			title = ""
		} else {
			doc := soup.HTMLParse(html)
			element := doc.Find("title")
			if element.Error != nil {
				//title = element.Error.Error()
				title = ""
			} else {
				title = element.Text()
			}
		}

		s.domainInfoCache.Store(domain, title)
	}

	return title.(string)
}
