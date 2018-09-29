package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric/common/tools/cryptogen/ca"
	"github.com/hyperledger/fabric/common/tools/cryptogen/msp"
)

type ActionAddOrderer struct {
	saveDir string // xxx/xxx/example.com
	*Args

	CA    *ca.CA
	TlsCA *ca.CA
}

func (self *ActionAddOrderer) Check(args *Args) error {
	self.Args = args

	if self.domain == "" {
		self.saveDir = args.dir
		self.domain = filepath.Base(args.dir)
	} else {
		self.saveDir = filepath.Join(args.dir, args.domain)
	}

	if self.name == "" {
		return errors.New("Need name")
	}

	self.CAPath = filepath.Join(self.saveDir, "ca")
	self.tlsCAPath = filepath.Join(self.saveDir, "tlsca")

	return nil
}

func (self *ActionAddOrderer) Run() (err error) {
	self.CA, err = LoadCA(self.CAPath)
	if err != nil {
		fmt.Println(self.CAPath)
		return err
	}

	self.TlsCA, err = LoadCA(self.tlsCAPath)
	if err != nil {
		return err
	}

	admin := fmt.Sprintf("%s@%s", "Admin", self.domain)
	ordererName := fmt.Sprintf("%s.%s", self.name, self.domain)
	adminCertPath := filepath.Join(self.saveDir, "users", admin, "msp", "signcerts", admin+"-cert.pem")

	generateNodes(filepath.Join(self.saveDir, "orderers"), ordererName, self.CA, self.TlsCA, msp.ORDERER, false)

	if err = copyAdminCert(adminCertPath, filepath.Join(self.saveDir, "msp", "admincerts")); err != nil {
		return err
	}
	os.Remove(filepath.Join(self.saveDir, "orderers", ordererName, "msp", "admincerts", ordererName+"-cert.pem"))

	if err = copyAdminCert(adminCertPath, filepath.Join(self.saveDir, "orderers", ordererName, "msp", "admincerts")); err != nil {
		return err
	}

	return nil
}
