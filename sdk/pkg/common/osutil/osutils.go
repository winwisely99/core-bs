package osutil

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
	"os"
)

// Exists returns whether the given file or directory exists or not.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func CheckAndMakeDir(path string) error {
	exists, err := Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

func RemoveDir(l *logger.Logger, path string) error {
	l.Infof("Removing directory....%s", path)
	_, err := SudoRunUnixCmd(true, `rm`, `-rf`, path)
	return err
}
