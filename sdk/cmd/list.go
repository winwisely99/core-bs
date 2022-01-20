package cmd

import (
	"fmt"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/colorutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/osutil"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/termutil"
	"github.com/spf13/cobra"
	"strings"
)

func NewListToolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all non-3rdparty tools available",
		Long:  "List all non-3rdparty tools available",
		RunE: func(cmd *cobra.Command, args []string) error {
			tools := []string{`lang`, `protofig`, `protoc-gen-validate`, `protoc-gen-config`, `googlesheetold`}
			contents := termutil.Contents{}
			for _, v := range tools {
				installed := osutil.BinExists(strings.TrimSpace(v))
				contents[v] = []string{colorutil.ColorRed(osutil.CrossMark) + " not installed"}
				if installed {
					contents[v] = []string{colorutil.ColorGreen(osutil.CheckMark) + " installed"}
				}
			}
			if _, err := fmt.Println(contents.String("All Tools")); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
