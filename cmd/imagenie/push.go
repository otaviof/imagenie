package main

import (
	"fmt"

	"github.com/otaviof/imagenie/pkg/imagenie"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pushCmd cobra definition for push sub-command
var pushCmd = &cobra.Command{
	Use:          "push <image> [image]",
	Short:        "Upload a container image to its registry.",
	PreRun:       setLogLevelCmd,
	RunE:         runPushCmd,
	SilenceUsage: true,
	Long: `### imagenie push

Upload container registry into container registry.
	`,
}

// init register sub-command in root.
func init() {
	flags := pushCmd.PersistentFlags()
	if err := viper.BindPFlags(flags); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(pushCmd)
}

// runPushCmd execute image push, more than one image is accepted
func runPushCmd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("a container image must be informed, as last argument")
	}

	for _, image := range args {
		m, err := imagenie.NewManager(image, image)
		if err != nil {
			return err
		}
		if err = m.Push(); err != nil {
			return err
		}
		if err = m.Delete(); err != nil {
			return err
		}
	}
	return nil
}
