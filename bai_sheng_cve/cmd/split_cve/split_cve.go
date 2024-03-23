package main

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/wuyyyyyou/go-share/pd"
)

func main() {
	df := pd.NewDataFrame()
	newdf := pd.NewDataFrame("更新")
	err := df.ReadExcel("../file/百胜-3.13.xlsx")
	if err != nil {
		logrus.Fatal(err)
	}
	newdf.SetHeads(df.GetHeads())
	rows := df.GetRows()
	var newRows [][]string
	for _, row := range rows {
		cveStr := row[6]
		cves := strings.Split(cveStr, "\n")
		for _, cve := range cves {
			newRow := make([]string, 6)
			copy(newRow, row[:6])
			newRow = append(newRow, cve)
			newRows = append(newRows, newRow)
		}
	}
	newdf.SetRows(newRows)
	err = newdf.SaveExcel("../file/百胜-3.13-new.xlsx")
}
