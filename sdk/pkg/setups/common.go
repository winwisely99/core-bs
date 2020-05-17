package setups

import (
	"bytes"
	"context"
	"errors"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/gitutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/pkgutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/oses"
	"html/template"
	"io/ioutil"
	"path/filepath"
)

type shellEnv struct {
	GoRoot               string
	GoPath               string
	FlutterPath          string
	DartPath             string
	JavaHome             string
	AndroidSdkHome       string
	AndroidNdkHome       string
	AndroidPlatformTools string
	AndroidTools         string
}

// Bootstrapper interface provides shared / common interface for
// bootstrapping OS setups.
type Bootstrapper interface {
	InstallOSPrequisites() error
	ShellEnv() error
}

var (
	foods = []string{
		"golang",
		"minikube",
		"kubectl",
		"hugo",
		"istioctl",
		"skaffold",
	}
)

func NewBootstrapper(c context.Context) (Bootstrapper, error) {
	u, err := oses.InitUserOsEnv()
	if err != nil {
		return nil, err
	}
	p, err := pkgutil.NewPkgUtil(c, u)
	if err != nil {
		return nil, err
	}
	return getBootstrapper(p)
}

func getBootstrapper(p *pkgutil.PkgUtil) (Bootstrapper, error) {
	switch p.GetPmName() {
	case "scoop":
		return &WinBootstrap{p}, nil
	case "dnf":
		return &RhelBootstrap{p}, nil
	case "apt":
		return &DebBootstrap{p}, nil
	case "brew":
		return &MacBootstrap{p}, nil
	default:
		return nil, errors.New("boostrapper error: unsupported os")
	}
}

func flatpaks(homepath string) (err error) {
	if err = osutil.CheckAndMakeDir(filepath.Join(homepath, "bin")); err != nil {
		return err
	}
	_, err = osutil.RunCmd(
		true,
		`flatpak`,
		`remote-add`,
		`--if-not-exists`,
		`flathub`,
		`https://flathub.org/repo/flathub.flatpakrepo`,
	)
	if err != nil {
		return err
	}
	if err = flatpakInstall("com.google.AndroidStudio"); err != nil {
		return err
	}
	if err = flatpakInstall("com.google.AndroidStudio"); err != nil {
		return err
	}
	astudio := `#!/usr/bin/env bash
	flatpak run coom.google.AndroidStudio
`
	return ioutil.WriteFile(filepath.Join(homepath, "bin", "astudio"), []byte(astudio), 0755)
}

func flatpakInstall(pkg string) error {
	_, err := osutil.RunCmd(true, `flatpak`,
		`install`, `flathub`, pkg)
	return err
}

func setupFlutter(l *logger.Logger, homepath string) error {
	flutterDir := filepath.Join(homepath, "flutter")
	err := gitutil.GitClone(l,
		"https://github.com/flutter/flutter.git",
		flutterDir,
	)
	if err != nil {
		return err
	}
	if err = gitutil.GitCheckout(l, flutterDir, "beta"); err != nil {
		return err
	}
	_, err = osutil.RunCmd(
		true,
		filepath.Join(flutterDir, "bin", `flutter`),
		`--enable-web`,
	)
	if err != nil { return err }
	_, err = osutil.RunCmd(
		true,
		filepath.Join(flutterDir, "bin", "cache", "dart-sdk", "bin", `pub`),
		`global`,
		`activate`, `protoc_plugin`,
	)
	return err
}

var (
	unixProfileTpl = `#!/usr/bin/env bash
GO111MODULE=on
GOROOT={{ .GoRoot }}
GOPATH={{ .GoPath }}
FLUTTER_PATH={{ .FlutterPath }}
DART_PATH={{ .DartPath }}
DART_HOMEPATH={{ .DartPath }}
JAVA_HOME={{ .JavaHome }}
ANDROID_SDK={{ .AndroidSdkHome }}
ANDROID_HOME={{ .AndroidSdkHome }}
ANDROID_NDK={{ .AndroidNdkHome }}
ANDROID_PLATFORM_TOOLS={{ .AndroidPlatformTools }}
ANDROID_TOOLS={{ .AndroidTools }}

export GOROOT GO111MODULE GOPATH JAVA_HOME ANDROID_HOME ANDROID_SDK ANDROID_NDK ANDROID_PLATFORM_TOOLS ANDROID_TOOLS
export PATH=$HOME/bin:$JAVA_HOME/bin:$DART_HOMEPATH:$GOPATH/bin:$FLUTTER_PATH:$DART_PATH:$GOROOT/bin:$PATH
`
)

func (s *shellEnv) execTemplate(b *bytes.Buffer) error{
	t, err := template.New("bs-profile.sh").Parse(unixProfileTpl)
	if err != nil {
		return err
	}
	return t.Execute(b, t)
}

func nixWriteProfile(homepath string, senv *shellEnv) error {
	b := bytes.NewBuffer(nil)
	if err := senv.execTemplate(b); err != nil {
		return err
	}
	profilePath := filepath.Join(homepath, ".bs-profile.sh")
	return ioutil.WriteFile(profilePath, b.Bytes(), 0755)
}