package gxutil

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	sh "github.com/ipfs/go-ipfs-api"
	manet "github.com/jbenet/go-multiaddr-net"
	ma "github.com/jbenet/go-multiaddr-net/Godeps/_workspace/src/github.com/jbenet/go-multiaddr"
)

func NewShell() *sh.Shell {
	ash, err := getLocalApiShell()
	if err == nil {
		return ash
	}

	return sh.NewShell("ipfs.io")
}

func getLocalApiShell() (*sh.Shell, error) {
	ipath := os.Getenv("IPFS_PATH")
	if ipath == "" {
		home := os.Getenv("HOME")
		if home == "" {
			return nil, errors.New("neither IPFS_PATH nor home dir set")
		}

		ipath = filepath.Join(home, ".ipfs")
	}

	apifile := filepath.Join(ipath, "api")

	data, err := ioutil.ReadFile(apifile)
	if err != nil {
		return nil, err
	}

	addr := strings.Trim(string(data), "\n\t ")

	maddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		return nil, err
	}

	_, host, err := manet.DialArgs(maddr)
	if err != nil {
		return nil, err
	}

	return sh.NewShell(host), nil
}