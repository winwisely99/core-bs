package setups

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/pkgutil"
	"path/filepath"
)

type MacBootstrap struct {
	*pkgutil.PkgUtil
}

func (m *MacBootstrap) ShellEnv() error {
	flutterPath := filepath.Join(m.Env.GetOsProperties().GetRoot(), "flutter", "bin")
	dartPath := filepath.Join(flutterPath, "cache", "dart-sdk", "bin")
	androidSdkHome := filepath.Join(m.Env.GetOsProperties().GetRoot(), "Library", "Android", "sdk")
	senv := &shellEnv{
		GoRoot: m.Env.GetGoRoot().Path(),
		GoPath: m.Env.GetGoPath().Path(),
		FlutterPath: flutterPath,
		DartPath: dartPath,
		JavaHome: "/Library/Java/JavaVirtualMachines/adoptopenjdk-12.0.2.jdk/Contents/Home",
		AndroidSdkHome: androidSdkHome,
		AndroidNdkHome: filepath.Join(androidSdkHome, "ndk-bundle"),
		AndroidPlatformTools: filepath.Join(androidSdkHome, "platform-tools"),
		AndroidTools: filepath.Join(androidSdkHome, "tools"),
	}
	return nixWriteProfile(m.Env.GetOsProperties().GetRoot(), senv)
}

func (m *MacBootstrap) InstallOSPrequisites() error {
	var err error
	homepath := m.Env.GetOsProperties().GetRoot()
	pm := m.GetPM()
	pkgs := []string{
		"git",
		"hyperkit",
		"docker-machine-hyperkit",
		"protobuf",
		"coreutils",
		"gcc",
		"libimobiledevice",
		"ideviceinstaller",
		"ios-deploy",
		"cocoapods",
	}
	if err = pm.Setup(); err != nil {
		return err
	}
	if err = pm.Install(pkgs...); err != nil {
		return err
	}
	if err = m.InstallFish(foods...); err != nil {
		return err
	}
	if err = m.caskInstall("adoptopenjdk8", "android-studio"); err != nil {
		return err
	}
	return setupFlutter(m.Log(), homepath)
}

func (m *MacBootstrap) caskInstall(pkgs ...string) error {
	cmds := append([]string{"cask", "install"}, pkgs...)
	_, err := osutil.RunCmd(
		true,
		`brew`,
		cmds...
	)
	return err
}
