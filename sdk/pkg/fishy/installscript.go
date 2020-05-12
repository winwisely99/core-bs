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
	} else if platform == "x86_64" || platform == "amd64" {
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

type installationDirs struct {
	installDir       string
	dlDir            string
	dlFilePath       string
	dlURL            string
	extractedDir     string
	extractedBinPath string
}

func (g *GoFishInstallation) determineDirs(arch string) *installationDirs {
	dlUrl := fmt.Sprintf("https://gofi.sh/releases/gofish-%s-%s-%s.tar.gz",
		g.Version, g.OSName, arch)
	dlDir := g.tempDir
	extractedDir := fmt.Sprintf(
		`%s%s%s-%s`, dlDir, g.separator, g.OSName, arch)
	installDir := fmt.Sprintf("%s%s%s",
		g.BinPath, g.separator, g.BinName)
	tarFile := fmt.Sprintf("%s%s%s.tar.gz",
		dlDir, g.separator, g.PkgName)
	extractedBinPath := fmt.Sprintf("%s%s%s",
		extractedDir, g.separator, g.BinName)
	return &installationDirs{
		installDir:       installDir,
		dlDir:            dlDir,
		dlFilePath:       tarFile,
		dlURL:            dlUrl,
		extractedDir:     extractedDir,
		extractedBinPath: extractedBinPath,
	}
}

func (g *GoFishInstallation) downloadBinary(d *installationDirs) error {
	c := fetcher.NewClient(g.l)
	res, err := c.Fetch(d.dlURL, "GET", nil, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	f, err := os.Create(d.dlFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, res.Body)
	return err
}

func (g *GoFishInstallation) installFile(d *installationDirs) error {
	if err := os.Chdir(d.dlDir); err != nil {
		return err
	}
	file, err := os.Open(d.dlFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	g.l.Debugf("extracting downloaded file: %s", d.dlFilePath)
	if err = osutil.ExtractTarGz(file); err != nil {
		return err
	}
	g.l.Debugf("reading binary file: %s", d.extractedBinPath)
	f, err := ioutil.ReadFile(d.extractedBinPath)
	if err != nil {
		log.Errorf("Cannot open path: %s: %v\n", d.extractedBinPath,
			err)
		return err
	}
	g.l.Infof("Installing to %s\n", g.BinPath)
	return ioutil.WriteFile(d.installDir, f, 0755)
}

func (g *GoFishInstallation) cleanDownloadedDirs(d *installationDirs) error {
	if err := os.Remove(d.dlFilePath); err != nil {
		return err
	}
	return os.RemoveAll(d.extractedDir)
}

func (g *GoFishInstallation) runInstallScript() error {
	arch := g.getArch()
	if !g.verifySupported() {
		return errors.New("architecture is not supported")
	}
	dirs := g.determineDirs(arch)
	if err := g.downloadBinary(dirs); err != nil {
		return err
	}
	if err := g.installFile(dirs); err != nil {
		return err
	}
	return g.cleanDownloadedDirs(dirs)
}
