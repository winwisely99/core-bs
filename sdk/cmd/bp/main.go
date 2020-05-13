package main

import (
	"flag"
	"fmt"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/gitutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	"github.com/sirupsen/logrus"
)

const (
	rtBoilGit = "https://github.com/getcouragenow/core-runtime"
)

var (
	tempDir = ""
	outDir  = ""
)

func main() {
	l := logger.NewLogger(logrus.DebugLevel, map[string]interface{}{
		"app": "gitstatik",
	})
	flag.StringVar(&tempDir, "t", "", "where the core-runtime will be cloned, it will be deleted once boilerplate is generated.")
	flag.StringVar(&outDir, "o", "", "where the embedded output will be generated")
	flag.Parse()
	if tempDir == "" {
		l.Fatal("Temporary dir has to be specified.")
	}
	if outDir == "" {
		l.Fatalf("Output dir has to be specified.")
	}
	sexists := osutil.BinExists("statik")
	if !sexists {
		if _, err := osutil.RunUnixCmd(true, `go`, `get`,
			`-u`, `-v`, `github.com/rakyll/statik`); err != nil {
			l.Fatalf("Error getting rakyll/statik: %v", err)
		}
	}
	statikBinPath, err := osutil.PrintBinaryPath("statik")
	if err != nil {
		l.Fatalf("Error getting statik bin path: %v", err)
	}
	l.Debugf("Creating temp directory: %s\n", tempDir)
	if err = osutil.CheckAndMakeDir(tempDir); err != nil {
		l.Fatalf("Cannot make temp directory")
	}
	l.Debugf("Creating output directory: %s\n", tempDir)
	if err = osutil.CheckAndMakeDir(outDir); err != nil {
		l.Fatalf("Cannot make output directory %s", outDir)
	}
	if err = gitutil.GitClone(l, rtBoilGit, tempDir); err != nil {
		l.Fatalf("Cannot clone %s to %s", rtBoilGit, tempDir)
	}
	l.Debugf("Embedding static files to binary in %s\n", outDir)
	_, err = osutil.RunUnixCmd(true, `statik`,
		fmt.Sprintf(`-src=%s/%s`, tempDir, "/boilerplate/"),
		`-ns`, `bp`, fmt.Sprintf(`-dest=%s`, outDir), `-f`)
	if err != nil {
		l.Fatalf("Cannot add statik file: %v", err)
	}
	l.Debug("Removing temporary dir")
	if _, err = osutil.RunUnixCmd(false, `rm`, `-rf`,
		tempDir); err != nil {
		l.Fatalf("Cannot remove temporary dir %s: %v", tempDir, err)
	}
	l.Debug("Uninstalling statik binary")
	if _, err = osutil.RunUnixCmd(false, `go`,
		`clean`, `-i`, `github.com/rakyll/statik`); err != nil {
		l.Fatalf("Cannot remove statik bin %s: %v", statikBinPath, err)
	}
}
