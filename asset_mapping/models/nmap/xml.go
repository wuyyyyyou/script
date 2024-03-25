package nmap

import (
	"encoding/xml"
)

type NmapRun struct {
	XMLName xml.Name `xml:"nmaprun"`
	Hosts   []Host   `xml:"host"`
}

type Host struct {
	XMLName xml.Name `xml:"host"`
	Status  Status   `xml:"status"`
	Address Address  `xml:"address"`
	Ports   Ports    `xml:"ports"`
}

type Status struct {
	XMLName   xml.Name `xml:"status"`
	State     string   `xml:"state,attr"`
	Reason    string   `xml:"reason,attr"`
	ReasonTtl string   `xml:"reason_ttl,attr"`
}

type Address struct {
	XMLName  xml.Name `xml:"address"`
	Addr     string   `xml:"addr,attr"`
	AddrType string   `xml:"addrtype,attr"`
}

type Ports struct {
	XMLName xml.Name `xml:"ports"`
	Port    []Port   `xml:"port"`
}

type Port struct {
	XMLName  xml.Name `xml:"port"`
	Protocol string   `xml:"protocol,attr"`
	Portid   string   `xml:"portid,attr"`
	State    State    `xml:"state"`
	Service  Service  `xml:"service"`
}

type State struct {
	XMLName   xml.Name `xml:"state"`
	State     string   `xml:"state,attr"`
	Reason    string   `xml:"reason,attr"`
	ReasonTtl string   `xml:"reason_ttl,attr"`
}

type Service struct {
	XMLName   xml.Name `xml:"service"`
	Name      string   `xml:"name,attr"`
	Servicefp string   `xml:"servicefp,attr"`
	Method    string   `xml:"method,attr"`
	Conf      string   `xml:"conf,attr"`
	Product   string   `xml:"product,attr"`
}
