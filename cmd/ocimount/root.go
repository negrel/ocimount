package ocimount

import (
	"os"

	"github.com/containers/storage/pkg/reexec"
	"github.com/spf13/cobra"
)

func init() {
	pflag := rootCmd.PersistentFlags()
	setupStoreOptionsFlags(pflag)
	setupLogrusFlags(pflag)
}

var rootCmd = &cobra.Command{}

func Execute() {
	if reexec.Init() {
		return
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
