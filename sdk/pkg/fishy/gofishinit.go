package fishy

import (
	"fmt"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/gitutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/oses"
	"os"
	"strings"
)

const (
	bsFishes = "github.com/getcouragenow/core-fish"
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
	separator    string
	userDir      string
	tempDir      string
	osProperties oses.OSInfoGetter
	l            *logger.Logger
}

func NewGoFishInstall(l *logger.Logger, u *oses.UserOsEnv) *GoFishInstallation {
	pkgName := "gofish"
	orgName := "fishworks"
	gitRepo := fmt.Sprintf("github.com/%s/%s", orgName, pkgName)
	goPath := u.GetGoEnv().GoPath()
	osName := strings.ToLower(u.GetOsProperties().GetOsInfo().GetOsName())
	separator := setSeparator(osName)
	if goPath == "" {
		os.Setenv("GOPATH", setGoPath(osName, u.GetOsProperties().GetRoot()))
	}
	g := &GoFishInstallation{
		Platform:     u.GetOsProperties().GetOsInfo().GetPlatform(),
		PkgName:      pkgName,
		OrgName:      orgName,
		Repo:         gitRepo,
		BinPath:      goPath + separator + "bin",
		SrcPath:      goPath + separator + "gofish",
		FishRepo:     "https://github.com/getcouragenow/core-fish",
		Version:      "v0.11.0",
		OSName:       osName,
		userDir:      u.GetOsProperties().GetRoot(),
		osProperties: u.GetOsProperties().GetOsInfo(),
		separator:    separator,
		l:            l,
	}
	g = g.setDirs()
	return g
}

func setSeparator(osname string) string {
	switch osname {
	case "windows":
		return `\`
	default:
		return "/"
	}
}

func setGoPath(sep, userdir string) string {
	return fmt.Sprintf(`%s%s%s%s%s`, userdir, sep,
		`workspace`, sep, "go")
}

func (g *GoFishInstallation) setDirs() *GoFishInstallation {
	osName := strings.ToLower(g.osProperties.GetOsName())
	switch osName {
	case "windows":
		g.BinName = g.PkgName + ".exe"
		g.tempDir = g.userDir + `\AppData\Local`
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
	_, err := osutil.RunUnixCmd(true,
		`gofish`, `init`)
	return err
}

func (g *GoFishInstallation) SetFishRig() error {
	os.Setenv("GOFISH_RIGS", g.SrcPath)
	return os.Setenv("GOFISH_DEFAULT_RIG", fmt.Sprintf("%s/%s", g.SrcPath, "core-fish"))
}

func (g *GoFishInstallation) InitGoFish() error {
	if err := g.SetFishRig(); err != nil {
		return err
	}
	return gitutil.GitClone(g.l, g.FishRepo, g.SrcPath)
}

func (g *GoFishInstallation) UninstallGoFish() error {
	return g.cleanGoFishGit()
}

func (g *GoFishInstallation) cleanGoFishGit() error {
	gitutil.GitRemove(g.l, "/usr/local/gofish")
	gitutil.GitRemove(g.l, fmt.Sprintf("%s/.%s", g.userDir, g.BinName))
	gitutil.GitRemove(g.l, g.SrcPath)
	gitutil.GitRemove(g.l, g.BinPath+"/"+g.BinName)
	return nil
}
