package fishy

import (
	"context"
	"fmt"
	home "github.com/fishworks/gofish/pkg/home"
	ctx2 "github.com/getcouragenow/core-bs/sdk/pkg/common/ctx"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/oses"
	"path/filepath"
	"strings"
)

type GoFishInstallation struct {
	PkgName      string
	BinName      string
	OrgName      string
	Repo         string
	BinPath      string
	SrcPath      string
	FishRepo     string
	Platform     string
	Version      string
	OSName       string
	userDir      string
	tempDir      string
	osProperties oses.OSInfoGetter
	l            *logger.Logger
}

func NewGoFishInstall(ctx context.Context, u *oses.UserOsEnv) *GoFishInstallation {
	l := ctx2.GetLogger(ctx)
	pkgName := "gofish"
	orgName := "fishworks"
	gitRepo := fmt.Sprintf("github.com/%s/%s", orgName, pkgName)
	osName := strings.ToLower(u.GetOsProperties().GetOsInfo().GetOsName())
	gopath := u.GetGoPath()
	g := &GoFishInstallation{
		Platform:     u.GetOsProperties().GetOsInfo().GetPlatform(),
		PkgName:      pkgName,
		OrgName:      orgName,
		Repo:         gitRepo,
		BinPath:      gopath.Path("bin"),
		SrcPath:      gopath.Path("gofish"),
		FishRepo:     "https://github.com/getcouragenow/core-fish",
		Version:      "v0.11.0",
		OSName:       osName,
		userDir:      u.GetOsProperties().GetRoot(),
		osProperties: u.GetOsProperties().GetOsInfo(),
		l:            l,
	}
	g = g.setDirs()
	return g
}

func (g *GoFishInstallation) setDirs() *GoFishInstallation {
	osName := strings.ToLower(g.osProperties.GetOsName())
	switch osName {
	case "windows":
		g.BinName = g.PkgName + ".exe"
		g.tempDir = filepath.Join(g.userDir, `AppData\Local`)
	default:
		g.BinName = g.PkgName
		g.tempDir = "/tmp"
	}
	return g
}

func (g *GoFishInstallation) InstallGoFish() error {
	g.l.Debugf("Installing gofish to GOPATH dir")
	// clean it up first
	g.cleanGoFishGit()
	return g.runInstallScript()
}

func (g *GoFishInstallation) GofishInit() error {
	g.l.Debugf("Running gofish init")
	_, err := osutil.RunCmd(true, `gofish`, `init`)
	if err != nil {
		return err
	}
	// TODO: add core-fish
	//if _, err = osutil.RunCmd(true, `gofish`, `rig`,
	//	`add`, g.FishRepo); err != nil {
	//	return err
	//}
	return nil
}

func (g *GoFishInstallation) UninstallGoFish() error {
	return g.cleanGoFishGit()
}

func (g *GoFishInstallation) cleanGoFishGit() error {
	// $HOME/.gofish
	osutil.RemoveDir(g.l, filepath.Join(g.userDir, fmt.Sprintf(".%s", g.BinName)))
	// $GOPATH/bin/gofish
	osutil.RemoveDir(g.l, filepath.Join(g.BinPath, g.BinName))
	// /usr/local/gofish
	osutil.RemoveDir(g.l, filepath.Join(home.HomePrefix, g.BinName))
	return nil
}
