package pkgutil

/*
pkgutil package provides wrapper around gofish commands.
*/

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
)

type PkgUtil struct {
	l *logger.Logger
}

func NewPkgUtil(l *logger.Logger) (*PkgUtil, error) {
	_, err := osutil.Exists("gofish")
	if err != nil {
		return nil, err
	}
	return &PkgUtil{
		l,
	}, nil
}

func (p *PkgUtil) InstallFish(pkgs ...string) error {
	p.l.Debugf("Installing %s\n", pkgs)
	params := []string{`install`}
	for _, pkg := range pkgs {
		params = append(params, pkg)
	}
	_, err := osutil.RunUnixCmd(true,
		`gofish`, params...)
	return err
}

func (p *PkgUtil) UninstallFish(pkgs ...string) error {
	p.l.Debugf("Uninstalling %s\n", pkgs)
	params := []string{`uninstall`}
	for _, pkg := range pkgs {
		params = append(params, pkg)
	}
	_, err := osutil.RunUnixCmd(
		true, `gofish`,
		params...)
	return err
}

func (p *PkgUtil) SearchFish(pkg string) error {
	p.l.Debugf("Searching package: %s\n", pkg)
	_, err := osutil.RunUnixCmd(true, `gofish`, `search`, pkg)
	return err
}

func (p *PkgUtil) RigsFish(cmds ...string) error {
	p.l.Debugf("Rigs operation %s\n", cmds)
	var params []string
	for _, pkg := range cmds {
		params = append(params, pkg)
	}
	_, err := osutil.RunUnixCmd(
		true, `gofish`,
		params...)
	return err
}

func (p *PkgUtil) UpdateFish(pkg string) error {
	p.l.Debugf("Update package: %s\n", pkg)
	_, err := osutil.RunUnixCmd(true,
		`gofish`, `update`, pkg)
	return err
}
