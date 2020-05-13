package main

import (
	"context"
	c "github.com/getcouragenow/core-bs/sdk/cmd"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	appName     = "bs"
	rootCmd     *cobra.Command
	globalUsage = `
	bs is a bootstrapper utility for kickstarting development with getcouragenow
	`
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          appName,
		Short:        globalUsage,
		Long:         globalUsage,
		SilenceUsage: true,
	}
	cmd.AddCommand(
		c.NewOsInfoCmd(),
		c.NewInstallFishCmd(),
		c.NewListToolsCmd(),
		// c.NewInitBoilerplateCmd(), // it's not going to be used here.
		c.NewPkgCmd(),
		// newSelfUpgradeCmd(),
	)
	return cmd
}

func main() {
	ctx := context.Background()
	l := logger.NewLogger(log.DebugLevel, map[string]interface{}{
		"app": "bs",
	})
	newCtx := context.WithValue(ctx, "logger", l)
	rootCmd = newRootCmd()
	if err := rootCmd.ExecuteContext(newCtx); err != nil {
		l.Fatalf("bs error: %v", err)
		os.Exit(1)
	}
}
