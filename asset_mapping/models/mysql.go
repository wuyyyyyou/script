package models

type Seed struct {
	BaseModel
	SeedName string    `gorm:"size:256;uniqueIndex"`
	Domains  []*Domain `gorm:"many2many:seed_domains;"`
}

type IP struct {
	BaseModel
	IP    string  `gorm:"size:128;uniqueIndex"`
	Ports []*Port `gorm:"foreignKey:BelongsIP;references:ID"`
}

type Port struct {
	BaseModel
	Port     int
	Protocol string `gorm:"size:64"`
	Finger   string `gorm:"size:128"`

	BelongsIP uint `gorm:"index"`
}

type Domain struct {
	BaseModel
	Domain string `gorm:"size:256;uniqueIndex"`
	IPs    []*IP  `gorm:"many2many:domain_ips;"`
}

type DomainIP struct {
	DomainID uint `gorm:"column:domain_id"`
	IPID     uint `gorm:"column:ip_id"`
}

type SeedDomain struct {
	SeedID   uint `gorm:"column:seed_id"`
	DomainID uint `gorm:"column:domain_id"`
}
