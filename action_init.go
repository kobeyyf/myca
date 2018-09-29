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
	saveDir string // xxx/xxx/example.com
	*Args

	CA    *ca.CA
	TlsCA *ca.CA
}

//
func (self *ActionInit) Check(args *Args) error {
	self.Args = args

	if self.domain == "" {
		self.saveDir = args.dir
		self.domain = filepath.Base(args.dir)
	} else {
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
func (self *ActionInit) Run() (err error) {

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
	err = msp.GenerateVerifyingMSP(mspDir, self.CA, self.TlsCA, false)
	if err != nil {
		return err
	}

	// 新建Admin用户
	// Admin@example.com
	// users/Admin@example.com
	admin := fmt.Sprintf("%s@%s", "Admin", self.domain)
	generateNodes(filepath.Join(self.saveDir, "users"), admin, self.CA, self.TlsCA, msp.CLIENT, false)

	adminCertPath := filepath.Join(self.saveDir, "users", admin, "msp", "signcerts",
		admin+"-cert.pem")
	os.RemoveAll(filepath.Join(self.saveDir, "msp", "admincerts"))
	if err = copyAdminCert(adminCertPath, filepath.Join(self.saveDir, "msp", "admincerts")); err != nil {
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
