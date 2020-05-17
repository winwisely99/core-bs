package setups

import "github.com/getcouragenow/core-bs/sdk/pkg/common/pkgutil"

type WinBootstrap struct {
	*pkgutil.PkgUtil
}

func (w *WinBootstrap) ShellEnv() error {
	return nil
}

func (w *WinBootstrap) InstallOSPrequisites() error {
	var err error
	homepath := w.Env.GetOsProperties().GetRoot()
	pm := w.GetPM()
	pkgs := []string{
		"git",
		"protobuf",
		"coreutils",
		"gcc",
		"make",
		"openjdk",
		"android-studio",
	}
	if err = pm.Setup(); err != nil {
		return err
	}
	if err = pm.Install(pkgs...); err != nil {
		return err
	}
	if err = w.InstallFish(foods...); err != nil {
		return err
	}
	return setupFlutter(w.Log(), homepath)
}