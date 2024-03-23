package analyze_log

type Asset struct {
	AssetID       string        `json:"AssetID"`
	AssetIPInfos  []AssetIPInfo `json:"AssetIpInfo"`
	AssetInfos    []AssetInfo   `json:"AssetInfo"`
	AssetName     string        `json:"AssetName"`
	AssetTag      string        `json:"AssetTag"`
	AssetType     string        `json:"AssetType"`
	DeviceLevel   string        `json:"DeviceLevel"`
	FoundTypeList string        `json:"FoundTypeList"`
	FoundTypeTime string        `json:"FoundTypeTime"`
	IsAccess      StringType    `json:"IsAccess"`
	IspCode       string        `json:"IspCode"`
	Location      string        `json:"Location"`
	NetPosition   StringType    `json:"NetPosition"`
	NetworkUnit   string        `json:"NetworkUnit"`
	ObjectName    string        `json:"ObjectName"`
	OrgCode       string        `json:"OrgCode"`
	State         StringType    `json:"State"`
	SystemName    string        `json:"SystemName"`
	TaskID        string        `json:"TaskID"`
}

type AssetIPInfo struct {
	IPType    StringType `json:"IpType"`
	IP        string     `json:"Ip"`
	OpenPorts []OpenPort `json:"OpenPort"`
}

type OpenPort struct {
	Protocol string     `json:"protocol"`
	Port     StringType `json:"port"`
}

type AssetInfo struct {
	Brands              []Brand       `json:"Brand"`
	Language            string        `json:"Language"`
	Memorys             []Memory      `json:"Memory"`
	OpenSource          StringType    `json:"OpenSource"`
	HardwareEnvironment string        `json:"HardwareEnvironment"`
	CPUs                []CPU         `json:"CPU"`
	Informations        []Information `json:"Information"`
	SoftwareEnvironment string        `json:"SoftwareEnvironment"`
	Chips               []Chip        `json:"Chip"`
}

type Brand struct {
	Manufacturer string `json:"Manufacturer"`
	Region       string `json:"Region"`
	Name         string `json:"Name"`
}

type Memory struct {
	Manufacturer string `json:"Manufacturer"`
	Model        string `json:"Model"`
}

type CPU struct {
	Manufacturer string `json:"Manufacturer"`
	Model        string `json:"Model"`
}

type Information struct {
	Version string `json:"Version"`
	Model   string `json:"Model"`
}

type Chip struct {
	Manufacturer string `json:"Manufacturer"`
	Model        string `json:"Model"`
}

func (a *Asset) FillNil() {
	if a.AssetIPInfos == nil || len(a.AssetIPInfos) == 0 {
		a.AssetIPInfos = []AssetIPInfo{
			{OpenPorts: []OpenPort{{}}},
		}
	}
	if a.AssetInfos == nil || len(a.AssetInfos) == 0 {
		a.AssetInfos = []AssetInfo{
			{
				Brands:       []Brand{{}},
				Memorys:      []Memory{{}},
				CPUs:         []CPU{{}},
				Informations: []Information{{}},
				Chips:        []Chip{{}},
			},
		}
	}
	if a.AssetIPInfos[0].OpenPorts == nil || len(a.AssetIPInfos[0].OpenPorts) == 0 {
		a.AssetIPInfos[0].OpenPorts = []OpenPort{{}}
	}
	if a.AssetInfos[0].Brands == nil || len(a.AssetInfos[0].Brands) == 0 {
		a.AssetInfos[0].Brands = []Brand{{}}
	}
	if a.AssetInfos[0].Memorys == nil || len(a.AssetInfos[0].Memorys) == 0 {
		a.AssetInfos[0].Memorys = []Memory{{}}
	}
	if a.AssetInfos[0].CPUs == nil || len(a.AssetInfos[0].CPUs) == 0 {
		a.AssetInfos[0].CPUs = []CPU{{}}
	}
	if a.AssetInfos[0].Informations == nil || len(a.AssetInfos[0].Informations) == 0 {
		a.AssetInfos[0].Informations = []Information{{}}
	}
	if a.AssetInfos[0].Chips == nil || len(a.AssetInfos[0].Chips) == 0 {
		a.AssetInfos[0].Chips = []Chip{{}}
	}
}
