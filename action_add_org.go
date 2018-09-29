package main

import (
	"errors"
	"fmt"
	// "io"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric/common/tools/cryptogen/ca"
	// "github.com/hyperledger/fabric/common/tools/cryptogen/csp"
	"github.com/hyperledger/fabric/common/tools/cryptogen/msp"
)

type ActionAddOrg struct {
	saveDir string // xxx/xxx/example.com
	*Args

	CA    *ca.CA
	TlsCA *ca.CA
}

//
func (self *ActionAddOrg) Check(args *Args) error {
	self.Args = args
	if self.name == "" {
		return errors.New("need org name")
	}

	if self.domain == "" {
		self.saveDir = args.dir
		self.domain = filepath.Base(args.dir)
	} else {
		self.domain = self.name + "." + self.domain
		self.saveDir = filepath.Join(args.dir, args.domain)
	}

	if self.CAPath == "" {
		self.CAPath = filepath.Join(self.saveDir, "ca")
	}
	if self.tlsCAPath == "" {
		self.tlsCAPath = filepath.Join(self.saveDir, "tlsca")
	}

	return nil
}

//
func (self *ActionAddOrg) Run() (err error) {

	// 获取ca 证书
	// 使用目录下的ca
	self.CA, err = LoadCA(self.CAPath)
	if err != nil || self.CA == nil {
		// fmt.Println(err)
		self.CA, err = GenCA(self.CAPath, "ca."+self.domain, self.Args)
		if err != nil {
			return err
		}
	}

	self.TlsCA, err = LoadCA(self.tlsCAPath)
	if err != nil {
		self.TlsCA, err = GenCA(self.tlsCAPath, "tlsca."+self.domain, self.Args)
		if err != nil {
			return err
		}
	}

	mspDir := filepath.Join(self.saveDir, "msp")
	err = msp.GenerateVerifyingMSP(mspDir, self.CA, self.TlsCA, true)
	if err != nil {
		return err
	}

	// 新建Admin用户
	// Admin@example.com
	// users/Admin@example.com
	admin := fmt.Sprintf("%s@%s", "Admin", self.domain)
	generateNodes(filepath.Join(self.saveDir, "users"), admin, self.CA, self.TlsCA, msp.CLIENT, true)

	adminCertPath := filepath.Join(self.saveDir, "users", admin, "msp", "signcerts",
		admin+"-cert.pem")
	os.RemoveAll(filepath.Join(self.saveDir, "msp", "admincerts"))
	if err = copyAdminCert(adminCertPath, filepath.Join(self.saveDir, "msp", "admincerts")); err != nil {
		return err
	}
	return nil
}
