package oses

/*
Package oses is for getting os, user, and git information
*/

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/fishworks/gofish/pkg/lazypath"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/gitutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/termutil"
)

// osProperties is the current environment for user's OS
type osProperties struct {
	name   string        // os username
	root   string        // root is the rootdir or homedir of the current user
	groups []*user.Group // groups user belongs to
	osInfo OSInfoGetter  // Os properties
}

func initOSProperties() (*osProperties, error) {
	shellUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	gids, err := shellUser.GroupIds()
	if err != nil {
		return nil, err
	}
	var groups []*user.Group
	for _, gid := range gids {
		group, err := user.LookupGroupId(gid)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	osInfo, err := getOsInfoGetter()
	if err != nil {
		return nil, err
	}
	return &osProperties{
		name:   shellUser.Username,
		root:   shellUser.HomeDir,
		groups: groups,
		osInfo: osInfo,
	}, nil
}
func (o *osProperties) GetName() string    { return o.name }
func (o *osProperties) GetRoot() string    { return o.root }
func (o *osProperties) GetAccount() string { return o.name }
func (o *osProperties) GetGroups() []*user.Group {
	return o.groups
}
func (o *osProperties) GetOsInfo() OSInfoGetter { return o.osInfo }
func (o *osProperties) ToContent() termutil.Contents {
	ms := termutil.Contents{}
	var groups []string
	for i := 0; i < len(o.GetGroups()); i++ {
		groups = append(groups, o.GetGroups()[i].Name)
	}
	ms["User"] = []string{o.GetName()}
	ms["User Homedir"] = []string{o.GetRoot()}
	ms["User Groups"] = groups
	return ms
}

// gitConfig is the current environment for user's git configuration
// users of this tool must have their own
type gitConfig struct {
	name    string       // git config --global user.name
	root    string       // ex: github.com/winwiselyxx
	account string       // git config --global user.email
	osInfo  OSInfoGetter // just added it for the sake of ease
}

func initGitConfig() (*gitConfig, error) {
	userName, err := RunCmd("git", "config", "user.name")
	if err != nil {
		return nil, err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	orig := "not a git dir"
	if gitutil.IsGitDir(pwd) {
		root, err := gitutil.GitRemoteInfo(pwd)
		if err != nil {
			orig = "not a git dir"
		}
		orig = *root
	}

	account, err := RunCmd("git", "config", "user.email")
	if err != nil {
		return nil, err
	}
	return &gitConfig{
		name:    *userName,
		root:    orig,
		account: *account,
		osInfo:  nil,
	}, nil
}

func (g *gitConfig) GetName() string    { return g.name }
func (g *gitConfig) GetRoot() string    { return g.root }
func (g *gitConfig) GetAccount() string { return g.account }
func (g *gitConfig) ToContent() termutil.Contents {
	ms := termutil.Contents{}
	ms["Git Global User"] = []string{g.GetName()}
	ms["Git Global Email"] = []string{g.GetAccount()}
	ms["Git Current URL"] = []string{g.GetRoot()}
	return ms
}

type GoConfig struct {
	goRoot lazypath.LazyPath
	goPath lazypath.LazyPath
}

func initGoConfig() *GoConfig {
	root := lazypath.LazyPath{
		EnvironmentVariable: "GOROOT",
		DefaultFn:           runtime.GOROOT,
	}
	goPath := lazypath.LazyPath{
		EnvironmentVariable: "GOPATH",
		DefaultFn:           setGoPath,
	}
	return &GoConfig{
		root,
		goPath,
	}
}
func (g *GoConfig) GoRoot() lazypath.LazyPath { return g.goRoot }
func (g *GoConfig) GoPath() lazypath.LazyPath { return g.goPath }
func (g *GoConfig) ToContent() termutil.Contents {
	ms := termutil.Contents{}
	ms["GOROOT"] = []string{g.GoRoot().Path()}
	ms["GOPATH"] = []string{g.GoPath().Path()}
	return ms
}
func setGoPath() string {
	if gp := os.Getenv("GOPATH"); gp != "" {
		return gp
	}
	u, _ := user.Current()
	workpaceGo := filepath.Join(u.HomeDir, "workspace", "go")
	osutil.CheckAndMakeDir(workpaceGo)
	return workpaceGo
}

type UserOsEnv struct {
	osProperties *osProperties
	goEnv        *GoConfig
	gitUser      *gitConfig
}

func InitUserOsEnv() (*UserOsEnv, error) {
	osProp, err := initOSProperties()
	if err != nil {
		return nil, err
	}
	gitUser, err := initGitConfig()
	if err != nil {
		return nil, err
	}
	goenv := initGoConfig()
	return &UserOsEnv{
		osProperties: osProp,
		gitUser:      gitUser,
		goEnv:        goenv,
	}, nil
}

func (u *UserOsEnv) GetGoPath() lazypath.LazyPath   { return u.goEnv.GoPath() }
func (u *UserOsEnv) GetGitUser() *gitConfig         { return u.gitUser }
func (u *UserOsEnv) GetGoRoot() lazypath.LazyPath   { return u.goEnv.GoRoot() }
func (u *UserOsEnv) GetGoEnv() *GoConfig            { return u.goEnv }
func (u *UserOsEnv) GetOsProperties() *osProperties { return u.osProperties }

func (u *UserOsEnv) PrintUserOsEnv() error {
	if _, err := fmt.Println(u.GetOsProperties().GetOsInfo().ToContent().String("Os and User info")); err != nil {
		return err
	}
	if _, err := fmt.Println(u.GetOsProperties().ToContent().String("")); err != nil {
		return err
	}
	if _, err := fmt.Println(u.GetGitUser().ToContent().String("")); err != nil {
		return err
	}
	if _, err := fmt.Println(u.GetGoEnv().ToContent().String("")); err != nil {
		return err
	}
	return nil
}
