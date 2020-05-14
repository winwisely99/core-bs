package cmd

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/fishy"
	"github.com/getcouragenow/core-bs/sdk/pkg/oses"
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

func NewInstallFishCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fish",
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			l := getLoggerFromContext(cmd.Context())
			l = l.AddFields(map[string]interface{}{
				"cmd": "fish",
			})
			newUserInfo, err := oses.InitUserOsEnv()
			if err != nil {
				return err
			}
			gf := fishy.NewGoFishInstall(l, newUserInfo)
			if installOpt {
				if err = gf.InstallGoFish(); err != nil {
					return err
				}
				// Post Installation
				gf.GofishInit()
				return err
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
