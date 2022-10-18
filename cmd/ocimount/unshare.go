package ocimount

import (
	"errors"
	"os"
	"os/exec"

	"github.com/containers/storage/pkg/unshare"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(unshareCmd)
	pflag := unshareCmd.PersistentFlags()
	setupStoreOptionsFlags(pflag)
	setupLogrusFlags(pflag)
}

var unshareCmd = &cobra.Command{
	Use:   "unshare",
	Short: "Run a command in a modified user namespace.",
	Run: func(cmd *cobra.Command, args []string) {
		if isRootless := unshare.IsRootless(); !isRootless {
			logrus.Error("please use unshare with rootless")
			os.Exit(1)
		}

		err := runUnshare(args)
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	},
}

func runUnshare(args []string) error {
	// exec the specified command, if there is one
	if len(args) < 1 {
		logrus.Debug("no cmd specified, detecting $SHELL...")

		// try to exec the shell, if one's set
		shell, shellSet := os.LookupEnv("SHELL")
		if !shellSet {
			return errors.New("no command specified and no $SHELL specified")
		}

		logrus.Debug("$SHELL detected: ", shell)
		args = []string{shell}
	}

	logrus.Debug("entering modified user namespace...")
	unshare.MaybeReexecUsingUserNamespace(false)
	logrus.Debugf("modified user namespace successfully entered, executing %v in a modified user namespace...", args)
	defer logrus.Debugf("%v executed.", args)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
