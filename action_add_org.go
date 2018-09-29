package main

import (
	"errors"
	// "fmt"
	// "io"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric/common/tools/cryptogen/ca"
	// "github.com/hyperledger/fabric/common/tools/cryptogen/csp"
	"github.com/hyperledger/fabric/common/tools/cryptogen/msp"
)

type ActionAddOrg struct {
	*Args

	CA    *ca.CA
	TlsCA *ca.CA
}

//
func (self *ActionAddOrg) Check(args *Args) error {
	self.Args = args
	if self.Name == "" {
		return errors.New("need org name: --name")
	}
	return nil
}

//
func (self *ActionAddOrg) Run() (err error) {
	self.CA, err = LoadCA(self.caDir)
	if err != nil {
		return err
	}

	self.TlsCA, err = LoadCA(self.tlscaDir)
	if err != nil {
		return err
	}

	mspDir := filepath.Join(self.Savedir, self.Domain, "msp")
	err = msp.GenerateVerifyingMSP(mspDir, self.CA, self.TlsCA, true)
	if err != nil {
		return err
	}

	admin := "Admin@" + self.Domain
	generateNodes(self.usersDir, admin, self.CA, self.TlsCA, msp.CLIENT, true)

	adminCertPath := filepath.Join(self.usersDir, admin, "msp", "signcerts",
		admin+"-cert.pem")

	admincertDir := filepath.Join(mspDir, "admincerts")
	os.RemoveAll(admincertDir)
	if err = copyAdminCert(adminCertPath, admincertDir); err != nil {
		return err
	}
	return nil
}
