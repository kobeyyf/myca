package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric/common/tools/cryptogen/ca"
	"github.com/hyperledger/fabric/common/tools/cryptogen/csp"
	"github.com/hyperledger/fabric/common/tools/cryptogen/msp"
)

type ActionInit struct {
	*Args

	CA    *ca.CA
	TlsCA *ca.CA
}

//
func (self *ActionInit) Check(args *Args) error {
	self.Args = args
	return nil
}

//
func (self *ActionInit) Run() (err error) {

	self.CA, err = LoadCA(self.caDir)
	if err != nil || self.CA == nil {
		// fmt.Println(err)
		self.CA, err = GenCA(self.caDir, "ca."+self.Domain, self.Args)
		if err != nil {
			return err
		}
	}

	self.TlsCA, err = LoadCA(self.tlscaDir)
	if err != nil || self.CA == nil {
		self.TlsCA, err = GenCA(self.tlscaDir, "tlsca."+self.Domain, self.Args)
		if err != nil {
			return err
		}
	}

	mspDir := filepath.Join(self.Savedir, self.Domain, "msp")
	err = msp.GenerateVerifyingMSP(mspDir, self.CA, self.TlsCA, false)
	if err != nil {
		return err
	}

	adminUser := "Admin@" + self.Domain
	generateNodes(self.usersDir, adminUser, self.CA, self.TlsCA, msp.CLIENT, false)

	adminUserCertPath := filepath.Join(self.usersDir, adminUser, "msp", "signcerts", adminUser+"-cert.pem")

	if err = copyAdminCert(adminUserCertPath, filepath.Join(mspDir, "admincerts")); err != nil {
		return err
	}
	return nil
}

func copyAdminCert(adminCertPath string, toDir string) error {
	name := filepath.Base(adminCertPath)
	_, err := os.Stat(toDir)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(toDir, 0755)
		} else {
			return err
		}
	}

	return copyFile(adminCertPath, filepath.Join(toDir, name))
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}

func generateNodes(baseDir string, nodeName string, signCA *ca.CA, tlsCA *ca.CA, nodeType int, nodeOUs bool) {

	nodeDir := filepath.Join(baseDir, nodeName)
	if _, err := os.Stat(nodeDir); os.IsNotExist(err) {
		err := msp.GenerateLocalMSP(nodeDir, nodeName, nil, signCA, tlsCA, nodeType, nodeOUs)
		if err != nil {

			fmt.Printf("Error generating local MSP for %s:\n%v\n", nodeName, err)
			os.Exit(1)
		}
	}
}

func LoadCA(dir string) (*ca.CA, error) {
	cert, err := ca.LoadCertificateECDSA(dir)
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, errors.New("LoadCertFailed " + dir)
	}

	_, key, err := csp.LoadPrivateKey(dir)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, errors.New("LoadKeyFailed")
	}

	return &ca.CA{
		Name:     cert.Subject.CommonName,
		Signer:   key,
		SignCert: cert,
	}, nil

}
func GenCA(dir string, domain string, args *Args) (*ca.CA, error) {
	os.RemoveAll(dir)
	return ca.NewCA(dir, args.Organizational, domain, args.Country, args.Province, args.Locality, args.OrganizationalUnit, args.StreetAddress, args.PostalCode)
}
