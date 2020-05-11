package pkgutil

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/fishy"
)

func InstallFish(pkg string, path string) error {
	if err := fishy.SetFishRig(path); err != nil {
		return err
	}
	_, err := osutil.RunUnixCmd(`gofish`, `install`, pkg)
	return err
}

func UninstallFish(pkg string, path string) error {
	if err := fishy.SetFishRig(path); err != nil {
		return err
	}
	_, err := osutil.RunUnixCmd(`gofish`, `uninstall`, pkg)
	return err
}

func SearchFish(pkg string, path string) error {
	if err := fishy.SetFishRig(path); err != nil {
		return err
	}
	_, err := osutil.RunUnixCmd(`gofish`, `search`, pkg)
	return err
}
