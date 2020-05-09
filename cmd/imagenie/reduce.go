package main

import (
	"fmt"
	"strings"

	"github.com/otaviof/imagenie/pkg/imagenie"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var reduceCmd = &cobra.Command{
	Use: "reduce",
	Run: runReduceCmd,
}

const pathSeparator = ":"

func init() {
	flags := reduceCmd.PersistentFlags()

	flags.StringSlice("copy", []string{}, "path to copy...")

	rootCmd.AddCommand(reduceCmd)
}

func getConfig(args []string) (*imagenie.Config, error) {
	if len(args) != 3 {
		exit(fmt.Errorf("not enough arguments %d", len(args)))
	}
	return &imagenie.Config{
		FromImage:   args[0],
		BaseImage:   args[1],
		TargetImage: args[2],
	}, nil
}

func prepareCopyPaths(copySlice []string) imagenie.CopyPaths {
	copyPaths := make(imagenie.CopyPaths, len(copySlice))
	for _, entry := range copyPaths {
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

func runReduceCmd(cmd *cobra.Command, args []string) {
	cfg, err := getConfig(args)
	if err != nil {
		exit(err)
	}

	i, err := imagenie.NewImagenie(cfg)
	if err != nil {
		exit(err)
	}

	copyPaths := prepareCopyPaths(viper.GetStringSlice("copy"))
	if err = i.Copy(copyPaths); err != nil {
		exit(err)
	}

	i.Labels()
}
