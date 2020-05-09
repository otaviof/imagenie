package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.Level(99))
}

// exit when on error display a final message and error.
func exit(err error) {
	fmt.Printf("[ERROR] %s\n", err)
	os.Exit(1)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
