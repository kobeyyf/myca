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
	*Args

	CA    *ca.CA
	TlsCA *ca.CA
}

func (self *ActionAddOrderer) Check(args *Args) error {
	self.Args = args

	if self.Name == "" {
		return errors.New("need peer name: --name")
	}
	return nil
}

func (self *ActionAddOrderer) Run() (err error) {
	self.CA, err = LoadCA(self.caDir)
	if err != nil {
		return err
	}

	self.TlsCA, err = LoadCA(self.tlscaDir)
	if err != nil {
		return err
	}

	ordererName := fmt.Sprintf("%s.%s", self.Name, self.Domain)
	ordererMsp := filepath.Join(self.orderersDir, ordererName, "msp")
	generateNodes(self.orderersDir, ordererName, self.CA, self.TlsCA, msp.ORDERER, false)
	os.Remove(filepath.Join(ordererMsp, "admincerts", ordererName+"-cert.pem"))

	adminUser := "Admin@" + self.Domain
	adminCertPath := filepath.Join(self.usersDir, adminUser, "msp", "signcerts", adminUser+"-cert.pem")

	if err = copyAdminCert(adminCertPath, filepath.Join(ordererMsp, "admincerts")); err != nil {
		return err
	}

	return nil
}
