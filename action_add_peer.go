package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric/common/tools/cryptogen/ca"
	"github.com/hyperledger/fabric/common/tools/cryptogen/msp"
)

type ActionAddPeer struct {
	*Args

	CA    *ca.CA
	TlsCA *ca.CA
}

func (self *ActionAddPeer) Check(args *Args) error {
	self.Args = args

	if self.Name == "" {
		return errors.New("Need peer name")
	}

	return nil
}

func (self *ActionAddPeer) Run() (err error) {
	self.CA, err = LoadCA(self.caDir)
	if err != nil {
		return err
	}

	self.TlsCA, err = LoadCA(self.tlscaDir)
	if err != nil {
		return err
	}

	peerName := fmt.Sprintf("%s.%s", self.Name, self.Domain)
	peerMsp := filepath.Join(self.peersDir, peerName, "msp")
	generateNodes(self.peersDir, peerName, self.CA, self.TlsCA, msp.PEER, true)

	admin := "Admin@" + self.Domain
	adminCertPath := filepath.Join(self.usersDir, admin, "msp", "signcerts", admin+"-cert.pem")
	if err = copyAdminCert(adminCertPath, filepath.Join(peerMsp, "admincerts")); err != nil {
		return err
	}

	os.Remove(filepath.Join(peerMsp, "admincerts", peerName+"-cert.pem"))

	if err = copyAdminCert(adminCertPath, filepath.Join(peerMsp, "admincerts")); err != nil {
		return err
	}

	return nil
}
