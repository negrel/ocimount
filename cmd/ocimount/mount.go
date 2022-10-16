package ocimount

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mountCmd)
	pflag := mountCmd.PersistentFlags()
	setupStoreOptionsFlags(pflag)
	setupLogrusFlags(pflag)
}

var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "Mount an OCI/Docker image.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 1 {
			return errors.New("expecting exactly one argument: an OCI/Docker image reference")
		}

		var mountpoint string
		mountpoint, err = mount(args[0])
		if err != nil {
			logrus.Errorf("failed to mount %q: %v", args[0], err)
			return nil
		}
		fmt.Println(mountpoint)

		return
	},
}

func mount(imgRefStr string) (mountpoint string, err error) {
	logrus.Debugf("mounting %q...", imgRefStr)

	store, err := containersStore()
	if err != nil {
		return
	}

	imgRef, err := parseReference(imgRefStr)
	if err != nil {
		return
	}

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		nbMount, merr := store.Mounted(imgRef.String())
		if merr != nil {
			logrus.Warnf("failed to check how many times %q is mounted: %v", imgRef, err)
		}
		logrus.Debugf("image already mounted %v time(s).", nbMount)
	}

	mountpoint, err = store.MountImage(imgRef.String(), []string{}, "")
	if err != nil {
		return
	}
	logrus.Infof("%q successfully mounted at %q.", imgRef, mountpoint)

	return
}
