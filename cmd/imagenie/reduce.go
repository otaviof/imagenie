package main

import (
	"fmt"
	"strings"

	"github.com/otaviof/imagenie/pkg/imagenie"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// copyFlag copy flag name.
	copyFlag = "copy"
	// pathSeparator copy flag separator
	pathSeparator = ":"
)

// reduceCmd cobra command definition.
var reduceCmd = &cobra.Command{
	Use:          "reduce <source-image> <base-image> <target-image> [options]",
	Short:        "Assemble target-image with parts of source-image on top of base-image",
	RunE:         runReduceCmd,
	PreRun:       setLogLevelCmd,
	SilenceUsage: true,
	Long: `### imagenie reduce

Assemble a new image using parts copied from "source-image" to "target-image", using "base-image"
base for newly created image. Metadata like image labels are also transported.

The objective of "reduce" command is to allow users to create lean-images out of arbitrary images,
and be a automation tool taking part of image building workflow.

Examples:

	Copying "/etc/os-release" and "/etc/alpine-release" to create "alpine:imagenie", storing files
	on "/tmp" directory.

	$ imagenie reduce \
		alpine:latest \     # source image
		alpine:latest \     # base image (runtime image)
		alpine:imagenie \   # target image
			--copy="/etc/os-release:/tmp" \
			--copy="/etc/alpine-release:/tmp" \
	`,
}

// init configure command-flags.
func init() {
	flags := reduceCmd.PersistentFlags()

	flags.StringSlice(copyFlag, []string{}, "copy data source-image '<source>:<destination>'")

	if err := viper.BindPFlags(flags); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(reduceCmd)
}

// getConfig expand arguments into a imagenie.Config instance.
func getConfig(args []string) (*imagenie.Config, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("not enough arguments %d", len(args))
	}
	return &imagenie.Config{
		FromImage:   args[0],
		BaseImage:   args[1],
		TargetImage: args[2],
	}, nil
}

// prepareCopyPaths intercept copy parameter to create a imagenie.CopyPaths map.
func prepareCopyPaths() imagenie.CopyPaths {
	copySlice := viper.GetStringSlice(copyFlag)
	copyPaths := make(imagenie.CopyPaths, len(copySlice))
	for _, entry := range copySlice {
		copyArgs := strings.Split(entry, pathSeparator)
		src := copyArgs[0]
		dst := ""
		if len(copyArgs) > 0 {
			dst = strings.Join(copyArgs[1:], pathSeparator)
		}
		copyPaths[src] = dst
	}
	return copyPaths
}

// runReduceCmd execute command.
func runReduceCmd(cmd *cobra.Command, args []string) error {
	cfg, err := getConfig(args)
	if err != nil {
		return err
	}

	i, err := imagenie.NewImagenie(cfg)
	if err != nil {
		return err
	}

	if err = i.Copy(prepareCopyPaths()); err != nil {
		return err
	}

	i.Labels()

	return nil
}
