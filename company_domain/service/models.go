package service

type DomainResponse struct {
	StateCode int              `json:"StateCode"`
	Results   []*DomainICPInfo `json:"Result"`
}

type DomainICPInfo struct {
	ID              uint `gorm:"primarykey" json:"-"`
	BelongToCompany uint `gorm:"index" json:"-"`

	// 站长之家返回的信息
	Domain      string `json:"Domain"`
	Owner       string `json:"Owner"`
	CompanyType string `json:"CompanyType"`
	SiteLicense string `json:"SiteLicense"`
	SiteName    string `json:"SiteName"`
	MainPage    string `json:"MainPage"`
	VerifyTime  string `json:"VerifyTime"`
}

type CompanyInfo struct {
	ID uint `gorm:"primarykey"`

	CompanyName string `gorm:"uniqueIndex;size:512"`
	Comment     string
	Type        string

	Credit string

	Finished bool `gorm:"default:false"`

	DomainInfos []*DomainICPInfo `gorm:"foreignKey:BelongToCompany;references:ID"`
}
