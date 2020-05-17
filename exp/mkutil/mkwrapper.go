package mkutil

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	log "github.com/sirupsen/logrus"
)

func runMake(args ...string) error {
	newArgs := []string{"-C"}
	for _, v := range args {
		newArgs = append(newArgs, v)
	}
	_, err := osutil.RunCmd(`make`, newArgs...)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func MakeHelpTool(path string) error {
	return runMake(path, `help`)
}

func MakePrintTool(path string) error {
	return runMake(path, `print`)
}

func MakeInstallTool(path string) error {
	return runMake(path, `build`)
}

func MakeTestTool(path string) error {
	return runMake(path, `test`)
}

func MakeUninstallTool(path string) error {
	return runMake(path, `build-clean`)
}

func MakeCustomCommand(path string, args ...string) error {
	newArgs := []string{path}
	for _, v := range args {
		newArgs = append(newArgs, v)
	}
	return runMake(newArgs...)
}
