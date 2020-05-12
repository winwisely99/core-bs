package fishy

import (
	"fmt"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/gitutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/oses"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	bsFishes = "github.com/getcouragenow/core-fish"
)

type GoFishInstallation struct {
	BinName  string
	OrgName  string
	Repo     string
	BinPath  string
	SrcPath  string
	FishRepo string
	Platform string
	Version  string
	OSName   string
	userDir  string
}

func NewGoFishInstall(u *oses.UserOsEnv) *GoFishInstallation {
	binName := "gofish"
	orgName := "fishworks"
	gitRepo := fmt.Sprintf("github.com/%s/%s", orgName, binName)
	goPath := u.GetGoEnv().GoPath()
	if goPath == "" {
		os.Setenv("GOPATH", fmt.Sprintf("%s/%s", u.GetOsProperties().GetRoot(), "workspace/go/"))
	}
	return &GoFishInstallation{
		Platform: u.GetOsProperties().GetOsInfo().GetPlatform(),
		BinName:  binName,
		OrgName:  orgName,
		Repo:     gitRepo,
		BinPath:  u.GetGoPath() + "/bin",
		SrcPath:  u.GetGoPath() + "/gofish",
		FishRepo: "https://github.com/getcouragenow/core-fish",
		Version:  "v0.11.0",
		OSName:   u.GetOsProperties().GetOsInfo().GetOsName(),
		userDir:  u.GetOsProperties().GetRoot(),
	}
}

func (g *GoFishInstallation) InstallGoFish() error {
	log.Infof("Installing gofish to GOPATH dir")
	// clean it up first
	g.cleanGoFishGit()
	return g.runInstallScript()
}

func (g *GoFishInstallation) GofishInit() error {
	_, err := osutil.RunUnixCmd(`gofish`, `init`)
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
	return gitutil.GitClone(g.FishRepo, g.SrcPath)
}

func (g *GoFishInstallation) UninstallGoFish() error {
	return g.cleanGoFishGit()
}

func (g *GoFishInstallation) cleanGoFishGit() error {
	os.RemoveAll("/usr/local/gofish")
	os.RemoveAll(fmt.Sprintf("%s/.%s", g.userDir, g.BinName))
	gitutil.GitRemove(g.SrcPath)
	gitutil.GitRemove(g.BinPath + "/" + g.BinName)
	return nil
}
