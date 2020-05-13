package cmd

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/oses"
	"github.com/spf13/cobra"
)

func NewOsInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "prints os info",
		RunE: func(cmd *cobra.Command, args []string) error {
			l := getLoggerFromContext(cmd.Context())
			l.AddFields(map[string]interface{}{
				"subcommand": "info",
			})
			newUserInfo, err := oses.InitUserOsEnv()
			if err != nil {
				return err
			}
			return newUserInfo.PrintUserOsEnv()
		},
	}
	return cmd
}
