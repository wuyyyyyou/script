package wangxingban

import (
	"github.com/wuyyyyyou/go-share/logutils"
)

var (
	Logger = (&logutils.LogFormatter{ReportCaller: true}).NewLogger()
)
