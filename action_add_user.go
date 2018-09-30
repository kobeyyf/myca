package main

import (
	"errors"
	"fmt"
	// "os"
	"path/filepath"

	"github.com/hyperledger/fabric/common/tools/cryptogen/ca"
	"github.com/hyperledger/fabric/common/tools/cryptogen/msp"
)

type ActionAddUser struct {
	*Args

	CA    *ca.CA
	TlsCA *ca.CA
}

func (self *ActionAddUser) Check(args *Args) error {
	self.Args = args

	if self.Name == "" {
		return errors.New("Need user name")
	}

	return nil
}

func (self *ActionAddUser) Run() (err error) {
	self.CA, err = LoadCA(self.caDir)
	if err != nil {
		return err
	}

	self.TlsCA, err = LoadCA(self.tlscaDir)
	if err != nil {
		return err
	}

	userName := fmt.Sprintf("%s@%s", self.Name, self.Domain)
	generateNodes(self.usersDir, userName, self.CA, self.TlsCA, msp.CLIENT, true)

	admin := "Admin@" + self.Domain
	adminCertPath := filepath.Join(self.usersDir, admin, "msp", "signcerts", admin+"-cert.pem")
	if err = copyAdminCert(adminCertPath, filepath.Join(self.usersDir, userName, "msp", "admincerts")); err != nil {
		return err
	}
	// os.Remove(filepath.Join(self.saveDir, "users", userName, "msp", "admincerts", userName+"-cert.pem"))

	// if err = copyAdminCert(adminCertPath, filepath.Join(self.saveDir, "users", userName, "msp", "admincerts")); err != nil {
	// 	return err
	// }

	return nil
}
