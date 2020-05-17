package pkgutil

/*
pkgutil package provides wrapper around gofish commands.
*/

import (
	"context"
	"errors"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/ctx"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/oses"
)

type PkgUtil struct {
	l *logger.Logger
	Env *oses.UserOsEnv
	pm PkgManager
}

func (p *PkgUtil) GetPmName() string {
	return p.pm.String()
}

func (p *PkgUtil) GetPM() PkgManager {
	return p.pm
}

func (p *PkgUtil) Log() *logger.Logger {
	return p.l
}

func NewPkgUtil(c context.Context, u *oses.UserOsEnv) (*PkgUtil, error) {
	var err error
	l := ctx.GetLogger(c)
	if u == nil {
		u, err = oses.InitUserOsEnv()
	}
	if err != nil {
		return nil, err
	}
	_, err = osutil.Exists("gofish")
	if err != nil {
		return nil, err
	}
	pm, err := getPkgManager(u.GetOsProperties().GetOsInfo().GetOsName())
	if err != nil {
		return nil, err
	}
	return &PkgUtil{
		l,
		u,
		pm,
	}, nil
}

func getPkgManager(osname string) (PkgManager, error) {
	switch osname {
	case "Windows":
		return &scoop{}, nil
	case "Darwin":
		return &brew{}, nil
	case "Linux":
		if osutil.BinExists("apt") {
			return &apt{}, nil
		}
		if osutil.BinExists("dnf") {
			return &dnf{}, nil
		}
	}
	return nil, errors.New("boostrapper error: unsupported OS")
}

func (p *PkgUtil) InstallFish(pkgs ...string) error {
	p.l.Debugf("Installing %s\n", pkgs)
	params := []string{`install`}
	for _, pkg := range pkgs {
		params = append(params, pkg)
	}
	_, err := osutil.RunCmd(true,
		`gofish`, params...)
	return err
}

func (p *PkgUtil) UninstallFish(pkgs ...string) error {
	p.l.Debugf("Uninstalling %s\n", pkgs)
	params := []string{`uninstall`}
	for _, pkg := range pkgs {
		params = append(params, pkg)
	}
	_, err := osutil.RunCmd(
		true, `gofish`,
		params...)
	return err
}

func (p *PkgUtil) SearchFish(pkg string) error {
	p.l.Debugf("Searching package: %s\n", pkg)
	_, err := osutil.RunCmd(true, `gofish`, `search`, pkg)
	return err
}

func (p *PkgUtil) RigsFish(cmds ...string) error {
	p.l.Debugf("Rigs operation %s\n", cmds)
	var params []string
	for _, pkg := range cmds {
		params = append(params, pkg)
	}
	_, err := osutil.RunCmd(
		true, `gofish`,
		params...)
	return err
}

func (p *PkgUtil) UpdateFish(pkg string) error {
	p.l.Debugf("Update package: %s\n", pkg)
	_, err := osutil.RunCmd(true,
		`gofish`, `update`, pkg)
	return err
}
