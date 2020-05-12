package main

import (
	"fmt"

	"github.com/otaviof/imagenie/pkg/imagenie"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pullCmd cobra definition for pull sub-command
var pullCmd = &cobra.Command{
	Use:          "pull <image> [image]",
	Short:        "Pull a upstream container image from registry.",
	PreRun:       setLogLevelCmd,
	RunE:         runPullCmd,
	SilenceUsage: true,
	Long: `### imagenie pull

Download a container image from upstream registry to local storage.
	`,
}

// init register sub-command in root.
func init() {
	flags := pullCmd.PersistentFlags()
	if err := viper.BindPFlags(flags); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(pullCmd)
}

// runPullCmd execute image pull.
func runPullCmd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("a container image must be informed, as last argument")
	}

	for _, image := range args {
		m, err := imagenie.NewManager(image, "")
		if err != nil {
			return err
		}
		if err = m.Pull(); err != nil {
			return err
		}
		if err = m.Delete(); err != nil {
			return err
		}
	}
	return nil
}
