package analyze_log

import (
	"fmt"
	"log"
	"testing"
)

func TestMain1(t *testing.T) {
	service := NewService()
	err := service.LoadAssets("files/content1.log")
	if err != nil {
		t.Fatal(err)
	}
	err = service.AssetsToXlsx("files/content01.xlsx")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain2(t *testing.T) {
	is := []int{1, 2, 3, 4}
	for _, i := range is {
		log.Printf("正在处理第 %v 个文件", i)
		src := fmt.Sprintf("./files/20231201/content_0%v.log", i)
		dst := fmt.Sprintf("./files/20231201/content_0%v.xlsx", i)

		service := NewService()
		err := service.LoadAssets(src)
		if err != nil {
			t.Fatal(err)
		}

		err = service.AssetsToXlsx(dst)
		if err != nil {
			t.Fatal(err)
		}
	}
}
