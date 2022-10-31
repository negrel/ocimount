package ocimount

import (
	"strings"

	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/storage"
	"github.com/containers/storage/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var storeOptions = types.StoreOptions{}

func setupStoreOptionsFlags(flagset *pflag.FlagSet) {
	flagset.StringVarP(&storeOptions.RunRoot, "run", "R", storeOptions.RunRoot, "Root of the runtime state tree")
	flagset.StringVarP(&storeOptions.GraphRoot, "graph", "g", storeOptions.GraphRoot, "Root of the storage tree")
	flagset.StringVarP(&storeOptions.GraphDriverName, "storage-driver", "s", storeOptions.GraphDriverName, "Storage driver to use")
}

func containersStore() (storage.Store, error) {
	var err error
	if storeOptions.GraphRoot == "" && storeOptions.RunRoot == "" &&
		storeOptions.GraphDriverName == "" && len(storeOptions.GraphDriverOptions) == 0 {
		storeOptions, err = types.DefaultStoreOptionsAutoDetectUID()
		if err != nil {
			return nil, err
		}
	}

	return storage.GetStore(storeOptions)
}

func parseReference(raw string) (ref reference.Reference, err error) {
	var namedRef reference.Named

	namedRef, err = reference.ParseNormalizedNamed(raw)
	switch {
	case err != nil:
		return
	case reference.IsNameOnly(namedRef):
		ref = reference.TagNameOnly(namedRef)
		if tagged, ok := ref.(reference.Tagged); ok {
			logrus.Infof("Using default tag: %s", tagged.Tag())
		}
	default:
		ref = namedRef
	}

	return
}

func setupLogrusFlags(flagset *pflag.FlagSet) {
	flagset.VarP(LogrusLevel{}, "level", "l", `Log level, one of "panic", "fatal", "error", "warn", "info", "debug", "trace"`)
}

type LogrusLevel struct{}

func (ll LogrusLevel) Set(v string) error {
	lvl, err := logrus.ParseLevel(strings.ToLower(v))
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)

	return nil
}

func (ll LogrusLevel) String() string {
	return logrus.GetLevel().String()
}

func (ll LogrusLevel) Type() string {
	return "log_level"
}
