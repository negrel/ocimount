package ocimount

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(umountCmd)
	pflag := umountCmd.PersistentFlags()
	setupStoreOptionsFlags(pflag)
	setupLogrusFlags(pflag)

	pflag.BoolP("force", "f", false, "force the unmount of an image")
}

var umountCmd = &cobra.Command{
	Use:   "umount",
	Short: "Unmount an OCI/Docker image.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("expecting exactly one argument: an OCI/Docker image reference")
		}

		flags := cmd.Flags()
		force, err := flags.GetBool("force")
		if err != nil {
			panic("force flag not found")
		}

		if err := umount(args[0], force); err != nil {
			logrus.Error("failed to unmount %q: %v", args[0], err)
			return nil
		}

		return nil
	},
}

func umount(imgRefStr string, force bool) (err error) {
	logrus.Debugf("unmounting %q...", imgRefStr)

	store, err := containersStore()
	if err != nil {
		return
	}

	imgRef, err := parseReference(imgRefStr)
	if err != nil {
		return
	}

	_, err = store.UnmountImage(imgRef.String(), force)
	if err != nil {
		return
	}
	logrus.Infof("%q successfully unmounted.", imgRef)

	return
}
