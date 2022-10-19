package ocimount

import (
	"errors"
	"os"

	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/storage"
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
	RunE:  runUmount,
}

func runUmount(cmd *cobra.Command, args []string) (err error) {
	// validating args
	if len(args) != 1 {
		return errors.New("expecting exactly one argument: an OCI/Docker image reference")
	}

	// getting flags
	flags := cmd.Flags()
	force, err := flags.GetBool("force")
	if err != nil {
		panic("force flag not found")
	}

	// parsing image arguments
	imgRef, err := parseReference(args[0])
	if err != nil {
		return err
	}

	// retrieving image store
	store, err := containersStore()
	if err != nil {
		logrus.Error("failed to retrieve container store: %v", err)
		return nil
	}

	logrus.Debugf("unmounting %q...", imgRef)
	if err := umount(store, imgRef, force); err != nil {
		logrus.Error("failed to unmount %q: %v", args[0], err)
		os.Exit(1)
	}
	logrus.Infof("%q successfully unmounted.", args[0])

	return nil
}

func umount(store storage.Store, imgRef reference.Reference, force bool) error {
	_, err := store.UnmountImage(imgRef.String(), force)
	return err
}
