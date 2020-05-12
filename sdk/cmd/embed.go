package cmd

import (
	"github.com/getcouragenow/core-bs/sdk/pkg/statics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	namespace string
)

func NewInitBoilerplateCmd() *cobra.Command {
	// joinedNS := strings.Join(statikNamespaces, "|")
	// usage := fmt.Sprintf("init -n [%s] <output_dir>", joinedNS)
	cmd := &cobra.Command{
		Use:   "init <output_dir>",
		Args:  cobra.ExactArgs(1),
		Short: "Write boilerplates to your specified directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println(namespace)
			bp, err := statics.NewBPAsset("bp")
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
