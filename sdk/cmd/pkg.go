package cmd

import (
	"context"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/ctx"
	"github.com/getcouragenow/core-bs/sdk/pkg/common/pkgutil"
	"github.com/spf13/cobra"
)

var (
	usage = "Install/Uninstall/Search 3rd party packages"
)

func NewPkgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "pkg [install|uninstall|search|rigs] <pkgname>",
		Short:        usage,
		Long:         usage,
		SilenceUsage: false,
	}

	cmd.AddCommand(
		installCmd(),
		uninstallCmd(),
		searchCmd(),
		rigsCmd(),
		updateCmd(),
	)
	return cmd
}

func installCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install <pkgname>",
		Short: "install 3rd party package",
		Long:  "install 3rd party package",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l := ctx.GetLogger(cmd.Context())
			l.AddFields(map[string]interface{}{
				"subcommand": "install",
			})
			newCtx := context.WithValue(cmd.Context(), "logger", l)
			p, err := pkgutil.NewPkgUtil(newCtx, nil)
			if err != nil {
				return err
			}
			return p.InstallFish(args...)
		},
	}
	return cmd
}

func uninstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall <pkgname>",
		Short: "uninstall 3rd party package",
		Long:  "uninstall 3rd party package",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l := ctx.GetLogger(cmd.Context())
			l.AddFields(map[string]interface{}{
				"subcommand": "uninstall",
			})
			newCtx := context.WithValue(cmd.Context(), "logger", l)
			p, err := pkgutil.NewPkgUtil(newCtx, nil)
			if err != nil {
				return err
			}
			return p.UninstallFish(args...)
		},
	}
	return cmd
}

func searchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <pkgname>",
		Short: "search 3rd party package",
		Long:  "search 3rd party package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l := ctx.GetLogger(cmd.Context())
			l.AddFields(map[string]interface{}{
				"subcommand": "search",
			})
			newCtx := context.WithValue(cmd.Context(), "logger", l)
			p, err := pkgutil.NewPkgUtil(newCtx, nil)
			if err != nil {
				return err
			}
			return p.SearchFish(args[0])
		},
	}
	return cmd
}

func rigsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rigs <operation>",
		Short: "rigs [list,add,remove,path]",
		Long:  "rigs [list,add,remove,path]",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l := ctx.GetLogger(cmd.Context())
			l.AddFields(map[string]interface{}{
				"subcommand": "rigs",
			})
			newCtx := context.WithValue(cmd.Context(), "logger", l)
			p, err := pkgutil.NewPkgUtil(newCtx, nil)
			if err != nil {
				return err
			}
			return p.RigsFish(args...)
		},
	}
	return cmd
}

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <pkgname>",
		Short: "update <pkgname>",
		Long:  "update <pkgname>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l := ctx.GetLogger(cmd.Context())
			l.AddFields(map[string]interface{}{
				"subcommand": "update",
			})
			newCtx := context.WithValue(cmd.Context(), "logger", l)
			p, err := pkgutil.NewPkgUtil(newCtx, nil)
			if err != nil {
				return err
			}
			return p.UpdateFish(args[0])
		},
	}
	return cmd
}
