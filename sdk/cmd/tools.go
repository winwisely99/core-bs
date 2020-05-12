package cmd

import (
	"fmt"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/colorutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/fishy"
	"github.com/getcouragenow/core-bs/sdk/pkg/oses"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
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
			newUserInfo, err := oses.InitUserOsEnv()
			if err != nil {
				return err
			}
			gf := fishy.NewGoFishInstall(newUserInfo)
			if installOpt {
				if err := gf.InstallGoFish(); err != nil {
					return err
				}
				// Post Installation
				gf.GofishInit()
				if err := gf.InitGoFish(); err != nil {
					return err
				}
				_, err = fmt.Fprintf(os.Stdout, fmt.Sprintf(`
				Please change GOFISH_RIGS environment variable to point to
					%s and GOFISH_DEFAULT_RIG to %s/%s
				in your shell profile.
			`, colorutil.ColorYellow(gf.SrcPath), colorutil.ColorCyan(gf.SrcPath), colorutil.ColorCyan("core-fish")))
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
