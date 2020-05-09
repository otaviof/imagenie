package main

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/otaviof/imagenie/pkg/imagenie"

	"github.com/spf13/cobra"
)

// logLevelFlag flag name.
const logLevelFlag = "log-level"

// rootCmd primary application cobra command.
var rootCmd = &cobra.Command{
	Use:    "imagenie <command>",
	Short:  "Utility tool to transform container images",
	PreRun: setLogLevelCmd,
}

// init instantiate flags.
func init() {
	flags := rootCmd.PersistentFlags()

	flags.Int(logLevelFlag, int(log.InfoLevel), "log verbosity level, from -2 to 3")

	if err := viper.BindPFlags(flags); err != nil {
		panic(err)
	}
}

// setLogLevelCmd use LOGLEVEL environment variable, and configure logrus.
func setLogLevelCmd(cmd *cobra.Command, args []string) {
	os.Setenv(imagenie.LogLevelEnv, strconv.Itoa(viper.GetInt(logLevelFlag)))
	imagenie.SetLogLevel()
}

func main() {
	imagenie.ReInit()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
