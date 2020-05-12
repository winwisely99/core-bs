package cmd

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/statics"
	"github.com/spf13/cobra"
)

func NewInitBoilerplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <output_dir>",
		Short: "Write boilerplates to your specified directory",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l := getLoggerFromContext(cmd.Context())
			l.Debugf("Writing boilerplate to...%s", args[0])
			bp, err := statics.NewBPAsset(l, "bp")
			if err != nil {
				return err
			}
			if err = bp.WriteAllFiles(args[0]); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
