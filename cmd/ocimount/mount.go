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

	// getting flags
	bind, err := cmd.Flags().GetString("bind")
	if err != nil {
		return
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

	// bind mount if necessary
	if bind != "" {
		logrus.Debug("binding %q to %q...", mountpoint, bind)

		err = mountBind(mountpoint, bind)
		if err != nil {
			logrus.Errorf("failed to bind mount %q to %q: %v", mountpoint, bind, err)
			logrus.Debug("cleaning up previous mount...")
			err = umount(store, imgRef, true)
			if err != nil {
				logrus.Debug("failed to clean up previous mount: %v", err)
			} else {
				logrus.Debug("previous mount successfully cleaned up.")
			}

			os.Exit(1)
		}

		logrus.Debugf("%q successfully bind mounted to %q.", imgRef, bind)
	}

	return

}

func mount(store storage.Store, imgRef reference.Reference) (mountpoint string, err error) {
	return store.MountImage(imgRef.String(), []string{}, "")
}

func mountBind(bind, to string) error {
	return stormount.Mount(bind, to, "", "rbind,rslave")
}

func umount(store storage.Store, imgRef reference.Reference, force bool) error {
	_, err := store.UnmountImage(imgRef.String(), force)
	return err
}
