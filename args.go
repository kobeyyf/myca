package main

import (
	"fmt"
	"os"
	"strings"
)

var (
	ERR_ARGS_TEMPLATE = "%s\n 这是模版"
)

type Args struct {
	command string

	dir    string // 证书所在目录
	domain string
	name   string // 名称

	// 使用自定义证书
	CAPath    string
	tlsCAPath string

	// csr模版
	Country            string
	Province           string
	Locality           string
	Organizational     string
	OrganizationalUnit string
	StreetAddress      string
	PostalCode         string

	EnabledNodeOUs bool
}

// O 组织
// OU 组织单位
// L 地址或城市
// ST 州（省）
// C 国家
func GetArgs() (*Args, error) {
	args := os.Args

	if len(args) < 2 {
		return nil, fmt.Errorf(ERR_ARGS_TEMPLATE, "缺少command")
	}

	a := &Args{
		command: args[1],
		dir:     "./",
		name:    "",
		domain:  "example.com",

		Country:            "CN",
		Province:           "FuJian",
		Locality:           "FuZhou",
		Organizational:     "ca",
		OrganizationalUnit: "myca",
		StreetAddress:      "",
		PostalCode:         "",
		EnabledNodeOUs:     true,
	}

	for _, arg := range args[2:] {
		fields := strings.SplitN(arg, "=", 2)
		if len(fields) != 2 {
			continue
		}
		switch fields[0] {
		case "--dir":
			a.dir = fields[1]
		case "--name":
			a.name = fields[1]
		case "--domain":
			a.domain = fields[1]
		case "--CAPath":
			a.CAPath = fields[1]
		case "--TlsCAPath":
			a.tlsCAPath = fields[1]
		case "--O":
			a.Organizational = fields[1]
		case "--OU":
			a.OrganizationalUnit = fields[1]
		case "--L":
			a.Province = fields[1]
		case "--ST":
			a.Locality = fields[1]
		case "--C":
			a.Country = fields[1]
		case "--EnabledNodeOUs":
			if fields[1] != "Y" && fields[1] != "y" {
				a.EnabledNodeOUs = false
			}
		}
	}
	return a, nil
}
