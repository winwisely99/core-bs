package fishy

import (
	"errors"
	"fmt"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/fetcher"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func (g *GoFishInstallation) getArch() string {
	platform := g.Platform
	if strings.Contains(platform, "arm") {
		return "arm"
	} else if strings.Contains(platform, "aarch64") {
		return "arm64"
	} else if platform == "x86" || platform == "i686" || platform == "386" {
		return "386"
	} else if platform == "x86_64" {
		return "amd64"
	}
	return ""
}

func (g *GoFishInstallation) verifySupported() bool {
	platform := g.getArch()
	if platform != "" {
		return true
	}
	return false
}

func (g *GoFishInstallation) downloadBinary(arch string, dldir string) error {
	dlUrl := fmt.Sprintf("https://gofi.sh/releases/gofish-%s-%s-%s.tar.gz",
		g.Version, strings.ToLower(g.OSName), arch)
	log.Infof("Downloading: %s\n", dlUrl)
	c := fetcher.NewClient()
	res, err := c.Fetch(dlUrl, "GET", nil, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	f, err := os.Create(fmt.Sprintf("%s/%s.tar.gz", dldir, g.BinName))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, res.Body)
	return err
}

func (g *GoFishInstallation) installFile(arch, dldir string) error {
	extractDir := fmt.Sprintf("%s/%s-%s", dldir, strings.ToLower(g.OSName), arch)
	installDir := fmt.Sprintf("%s/%s", g.BinPath, g.BinName)
	tarFile := fmt.Sprintf("%s/%s.tar.gz", dldir, g.BinName)
	if err := os.Chdir(dldir); err != nil {
		return err
	}
	file, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := osutil.ExtractTarGz(file); err != nil {
		return err
	}
	if err := os.Remove(tarFile); err != nil {
		return err
	}
	log.Infof("Installing to %s\n", extractDir)
	f, err := ioutil.ReadFile(fmt.Sprintf("%s/gofish", extractDir))
	if err != nil {
		log.Errorf("Cannot open path: %s: %v\n", extractDir, err)
		return err
	}
	if err := ioutil.WriteFile(installDir, f, 0755); err != nil {
		return err
	}
	return os.RemoveAll(extractDir)
}

func (g *GoFishInstallation) runInstallScript() error {
	arch := g.getArch()
	dldir := "/tmp"
	if !g.verifySupported() {
		return errors.New("architecture is not supported")
	}
	if err := g.downloadBinary(arch, dldir); err != nil {
		return err
	}
	return g.installFile(arch, dldir)
}
