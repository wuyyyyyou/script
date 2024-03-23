package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/samber/lo"

	"github.com/wuyyyyou/script/company_domain/utils"
)

// 读取excel文件，获取域名，生成Quake的查询语句
func TestTmp_GetQuakeQuery(t *testing.T) {
	t.Log(utils.GetQuakeQuery("/Users/leyouming/program/golang_note/golang_code/go_code/script/company_domain/file/测绘2.xlsx"))
}

func TestTmp_LookupIP(t *testing.T) {
	address := "1c5e7647.src.ucloud.com.cn"
	//address := "101.34.71.193"
	IPs, err := net.LookupIP(address)
	if err != nil {
		t.Fatal(err)
	}
	for _, ip := range IPs {
		t.Log(fmt.Sprintf("%v", ip))
	}
}

func TestTmp_GetIP(t *testing.T) {
	input := "/Users/leyouming/program/golang_note/golang_code/go_code/script/company_domain/file/all_domain.xlsx"
	output := "/Users/leyouming/program/golang_note/golang_code/go_code/script/company_domain/file/ips.xlsx"
	err := utils.GetIP(input, output)
	if err != nil {
		t.Fatal(err)
	}
}

// 分割ip.txt文件
func TestTmp_SplitTxt(t *testing.T) {
	filePath := "/Users/leyouming/company_program/scan_tool/AssetMapping/nmap/ips20240322.txt"

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)
	ips := strings.Split(contentStr, "\n")
	chunks := lo.Chunk(ips, 3)

	for i, chunk := range chunks {
		file, err := os.Create(strings.Replace(filePath, ".txt", fmt.Sprintf("_%d.txt", i), 1))
		if err != nil {
			fmt.Println("创建文件时出错:", err)
			return
		}
		writer := bufio.NewWriter(file)
		for _, line := range chunk {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				t.Fatal(err)
			}
		}

		if err := writer.Flush(); err != nil {
			t.Fatal(err)
		}

		_ = file.Close()
	}

}
