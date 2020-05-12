package pkgutil

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
)

func InstallFish(pkg string, path string) error {
	_, err := osutil.RunUnixCmd(true, `gofish`, `install`, pkg)
	return err
}

func UninstallFish(pkg string, path string) error {
	_, err := osutil.RunUnixCmd(true, `gofish`, `uninstall`, pkg)
	return err
}

func SearchFish(pkg string, path string) error {
	_, err := osutil.RunUnixCmd(true, `gofish`, `search`, pkg)
	return err
}
