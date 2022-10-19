package ocimount

import (
	"errors"
	"fmt"
	"os"

	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/storage"
	stormount "github.com/containers/storage/pkg/mount"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mountCmd)
	pflag := mountCmd.PersistentFlags()
	setupStoreOptionsFlags(pflag)
	setupLogrusFlags(pflag)

	pflag.StringP("bind", "b", "", "bind mount to this directory")
}

var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "Mount an OCI/Docker image.",
	RunE:  runMount,
}

func runMount(cmd *cobra.Command, args []string) (err error) {
	// validating args
	if len(args) != 1 {
		return errors.New("expecting exactly one argument: an OCI/Docker image reference")
	}

	// parsing image arguments
	imgRef, err := parseReference(args[0])
	if err != nil {
		return
	}

	// retrieving image store
	store, err := containersStore()
	if err != nil {
		logrus.Error("failed to retrieve container store: %v", err)
		return nil
	}

	// check if image is already mounted
	if logrus.IsLevelEnabled(logrus.InfoLevel) {
		nbMount, merr := store.Mounted(imgRef.String())
		if merr != nil {
			logrus.Warnf("failed to check how many times %q is mounted: %v", imgRef, err)
		}
		if nbMount > 0 {
			logrus.Infof("image already mounted %v time(s).", nbMount)
		}
	}

	// mounting the image as read only
	logrus.Debugf("mounting %q...", imgRef)
	var mountpoint string
	mountpoint, err = mount(store, imgRef)
	if err != nil {
		logrus.Errorf("failed to mount %q: %v", imgRef, err)
		os.Exit(1)
	}
	logrus.Infof("%q successfully mounted at %q.", imgRef, mountpoint)
	fmt.Println(mountpoint)

	return
}

func mount(store storage.Store, imgRef reference.Reference) (mountpoint string, err error) {
	return store.MountImage(imgRef.String(), []string{}, "")
}
