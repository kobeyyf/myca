package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ERR_ARGS_TEMPLATE = "%s\n 这是模版"
)

type Args struct {
	command string

	Savedir string // 证书所在目录
	Domain  string // 域名
	Name    string // 名称

	// csr模版
	Country            string
	Province           string
	Locality           string
	Organizational     string
	OrganizationalUnit string
	StreetAddress      string
	PostalCode         string

	// 路径信息
	caDir       string
	tlscaDir    string
	usersDir    string
	orderersDir string
	peersDir    string
}

// O 组织
// OU 组织单位
// L 地址或城市
// ST 州（省）
// C 国家
func GetArgs() (*Args, error) {
	args := os.Args

	if len(args) < 2 {
		return nil, fmt.Errorf(ERR_ARGS_TEMPLATE, "缺少 command")
	}

	a := &Args{
		command: args[1],
		Savedir: "./",
		Name:    "",
		Domain:  "",

		Country:            "CN",
		Province:           "FuJian",
		Locality:           "FuZhou",
		Organizational:     "ca",
		OrganizationalUnit: "myca",
		StreetAddress:      "",
		PostalCode:         "",
	}

	for _, arg := range args[2:] {
		fields := strings.SplitN(arg, "=", 2)
		if len(fields) != 2 {
			continue
		}
		switch fields[0] {
		case "--save_dir":
			a.Savedir = fields[1]
		case "--name":
			a.Name = fields[1]
		case "--domain":
			a.Domain = fields[1]
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
		}
	}

	if a.Domain == "" {
		return nil, fmt.Errorf(ERR_ARGS_TEMPLATE, "缺少域名 --domain")
	}
	a.caDir = filepath.Join(a.Savedir, a.Domain, "ca")
	a.usersDir = filepath.Join(a.Savedir, a.Domain, "users")
	a.tlscaDir = filepath.Join(a.Savedir, a.Domain, "tlsca")
	a.orderersDir = filepath.Join(a.Savedir, a.Domain, "orderers")
	a.peersDir = filepath.Join(a.Savedir, a.Domain, "peers")

	return a, nil
}
