package cmd

import (
	"context"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/ctx"
	"github.com/getcouragenow/core-bs/sdk/pkg/fishy"
	"github.com/getcouragenow/core-bs/sdk/pkg/oses"
	"github.com/getcouragenow/core-bs/sdk/pkg/setups"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	short = `install / uninstall gofish to your PATH`
	long  = `use 'bs fish -i' to install gofish 
	and 'bs fish -u' to uninstall gofish
	`
)

var (
	installOpt, uninstallOpt bool
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init -[i|u]",
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			l := ctx.GetLogger(cmd.Context())
			l = l.AddFields(map[string]interface{}{
				"cmd": "fish",
			})
			newCtx := context.WithValue(cmd.Context(), "logger", l)
			newUserInfo, err := oses.InitUserOsEnv()
			if err != nil {
				return err
			}
			gf := fishy.NewGoFishInstall(newCtx, newUserInfo)
			if installOpt {
				if err = gf.InstallGoFish(); err != nil {
					return err
				}
				// Post Installation
				if err = gf.GofishInit(); err != nil {
					return err
				}
				newSetup, err := setups.NewBootstrapper(newCtx)
				if err != nil { return err }
				if err = newSetup.ShellEnv(); err != nil { return err }
				return newSetup.InstallOSPrequisites()
			}
			if uninstallOpt {
				log.Infof("Uninstalling gofish binary and sources")
				return gf.UninstallGoFish()
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVarP(&installOpt, "install", "i",
		false, `use this flag to install gofish`)
	cmd.PersistentFlags().BoolVarP(&uninstallOpt, "uninstall", "u",
		false, `use this flag to uninstall gofish`)
	return cmd
}
