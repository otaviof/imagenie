package main

import (
	"fmt"
	"strings"

	"github.com/otaviof/imagenie/pkg/imagenie"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// cmdFlag cmd flag name.
	cmdFlag = "cmd"
	// copyFlag copy flag name.
	copyFlag = "copy"
	// entrypointFlag entrypoint flag name.
	entrypointFlag = "entrypoint"
	// pushFlag push flag name
	pushFlag = "push"
	// runFlag run flag name.
	runFlag = "run"
)

// pathSeparator copy flag separator
const pathSeparator = ":"

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

	flags.StringSlice(entrypointFlag, []string{}, "overwrite target-image 'entrypoint'")
	flags.StringSlice(cmdFlag, []string{}, "overwrite target-image 'cmd'")
	flags.StringSlice(copyFlag, []string{}, "copy data source-image '<source>:<destination>'")
	flags.StringSlice(runFlag, []string{}, "run arbitrary command on target-image")
	flags.Bool(pushFlag, false, "push target image to registry")

	if err := viper.BindPFlags(flags); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(reduceCmd)
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

// newImagenie instantiate with configuration based on command-line arguments.
func newImagenie(args []string) (*imagenie.Imagenie, error) {
	cfg, err := getConfig(args)
	if err != nil {
		return nil, err
	}
	return imagenie.NewImagenie(cfg)
}

// runReduceCmd execute command.
func runReduceCmd(cmd *cobra.Command, args []string) error {
	i, err := newImagenie(args)
	if err != nil {
		return err
	}

	// copying data, and copying labels
	if err = i.Copy(prepareCopyPaths()); err != nil {
		return err
	}
	i.Labels()

	// make settings in the image, entrypoint and cmd
	if entrypoint := viper.GetStringSlice(entrypointFlag); len(entrypoint) > 0 {
		i.TargetMgr.SetEntrypoint(entrypoint)
	}
	if cmd := viper.GetStringSlice(cmdFlag); len(cmd) > 0 {
		i.TargetMgr.SetCMD(cmd)
	}

	// run additional commands on target image
	if runCommands := viper.GetStringSlice(runFlag); len(runCommands) > 0 {
		i.RunAll(runCommands)
	}

	// commit target image, clean-up mounts and container
	if err = i.TargetMgr.Commit(); err != nil {
		return err
	}

	// uploading target container image
	if viper.GetBool(pushFlag) {
		if err = i.TargetMgr.Push(); err != nil {
			return err
		}
	}

	// final clean up, unmount and remove working containers
	return i.CleanUp()
}
