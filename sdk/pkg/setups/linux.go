package setups

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/common/pkgutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/oses"
	"path/filepath"
)

type DebBootstrap struct {
	*pkgutil.PkgUtil
}

type RhelBootstrap struct {
	*pkgutil.PkgUtil
}

func linuxWriteProfile(u *oses.UserOsEnv) error {
	flutterPath := filepath.Join(u.GetOsProperties().GetRoot(), "flutter", "bin")
	dartPath := filepath.Join(flutterPath, "cache", "dart-sdk", "bin")
	androidSdkHome := filepath.Join(u.GetOsProperties().GetRoot(), ".Android", "sdk")
	senv := &shellEnv{
		GoRoot: u.GetGoRoot().Path(),
		GoPath: u.GetGoPath().Path(),
		FlutterPath: flutterPath,
		DartPath: dartPath,
		JavaHome: "/usr/lib/jvm/java-8-openjdk-amd64",
		AndroidSdkHome: androidSdkHome,
		AndroidNdkHome: filepath.Join(androidSdkHome, "ndk-bundle"),
		AndroidPlatformTools: filepath.Join(androidSdkHome, "platform-tools"),
		AndroidTools: filepath.Join(androidSdkHome, "tools"),
	}
	return nixWriteProfile(u.GetOsProperties().GetRoot(), senv)
}

func (r *RhelBootstrap) ShellEnv() error {
	return linuxWriteProfile(r.Env)
}

func (d *DebBootstrap) ShellEnv() error {
	return linuxWriteProfile(d.Env)
}

func (d *DebBootstrap) InstallOSPrequisites() error {
	var err error
	homepath := d.Env.GetOsProperties().GetRoot()
	pm := d.GetPM()
	pkgs := []string{
		"libprotobuf-dev",
		"protobuf-compiler",
		"curl",
		"unzip",
		"git",
		"openssh-server",
		"gcc",
		"openjdk-8-jdk",
		"flatpak",
		"qemu-kvm",
		"libvirt-dev",
		"virt-manager",
		"build-essential",
	}
	if err = pm.Setup(); err != nil {
		return err
	}
	if err = pm.Install(pkgs...); err != nil {
		return err
	}
	if err = d.InstallFish(foods...); err != nil {
		return err
	}
	if err = flatpaks(homepath); err != nil {
		return err
	}
	return setupFlutter(d.Log(), homepath)
}

func (r *RhelBootstrap) InstallOSPrequisites() error {
	var err error
	homepath := r.Env.GetOsProperties().GetRoot()
	pm := r.GetPM()
	pkgs := []string{
		"curl",
		"unzip",
		"git",
		"openssh-server",
		"gcc",
		"java-8-openjdk",
		"qemu-kvm",
		"qemu-img",
		"virt-manager",
		"libvirt",
		"libvirt-python",
		"libvirt-client",
		"virt-install",
		"virt-viewer",
		"bridge-utils",
		"protobuf-devel",
	}
	if err = pm.Setup(); err != nil {
		return err
	}
	if err = pm.Install(pkgs...); err != nil {
		return err
	}
	if err = r.InstallFish(foods...); err != nil {
		return err
	}
	if err = flatpaks(homepath); err != nil {
		return err
	}
	return setupFlutter(r.Log(), homepath)
}
