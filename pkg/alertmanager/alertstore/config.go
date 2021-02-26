package alertstore

import (
	"context"
	"flag"
	"fmt"

	"github.com/pkg/errors"

	"github.com/cortexproject/cortex/pkg/alertmanager/alertstore/configdb"
	"github.com/cortexproject/cortex/pkg/alertmanager/alertstore/local"
	"github.com/cortexproject/cortex/pkg/alertmanager/alertstore/objectclient"
	"github.com/cortexproject/cortex/pkg/chunk"
	"github.com/cortexproject/cortex/pkg/chunk/aws"
	"github.com/cortexproject/cortex/pkg/chunk/azure"
	"github.com/cortexproject/cortex/pkg/chunk/gcp"
	"github.com/cortexproject/cortex/pkg/configs/client"
)

// LegacyConfig configures the alertmanager backend using the legacy storage clients.
type LegacyConfig struct {
	Type     string        `yaml:"type"`
	ConfigDB client.Config `yaml:"configdb"`

	// Object Storage Configs
	Azure azure.BlobStorageConfig `yaml:"azure"`
	GCS   gcp.GCSConfig           `yaml:"gcs"`
	S3    aws.S3Config            `yaml:"s3"`
	Local local.StoreConfig       `yaml:"local"`
}

// RegisterFlags registers flags.
func (cfg *LegacyConfig) RegisterFlags(f *flag.FlagSet) {
	cfg.ConfigDB.RegisterFlagsWithPrefix("alertmanager.", f)
	f.StringVar(&cfg.Type, "alertmanager.storage.type", "configdb", "Type of backend to use to store alertmanager configs. Supported values are: \"configdb\", \"gcs\", \"s3\", \"local\".")

	cfg.Azure.RegisterFlagsWithPrefix("alertmanager.storage.", f)
	cfg.GCS.RegisterFlagsWithPrefix("alertmanager.storage.", f)
	cfg.S3.RegisterFlagsWithPrefix("alertmanager.storage.", f)
	cfg.Local.RegisterFlags(f)
}

// Validate config and returns error on failure
func (cfg *LegacyConfig) Validate() error {
	if err := cfg.Azure.Validate(); err != nil {
		return errors.Wrap(err, "invalid Azure Storage config")
	}
	if err := cfg.S3.Validate(); err != nil {
		return errors.Wrap(err, "invalid S3 Storage config")
	}
	return nil
}

// NewLegacyAlertStore returns a new rule storage backend poller and store
func NewLegacyAlertStore(cfg LegacyConfig) (AlertStore, error) {
	switch cfg.Type {
	case "configdb":
		c, err := client.New(cfg.ConfigDB)
		if err != nil {
			return nil, err
		}
		return configdb.NewStore(c), nil
	case "azure":
		return newLegacyObjAlertStore(azure.NewBlobStorage(&cfg.Azure))
	case "gcs":
		return newLegacyObjAlertStore(gcp.NewGCSObjectClient(context.Background(), cfg.GCS))
	case "s3":
		return newLegacyObjAlertStore(aws.NewS3ObjectClient(cfg.S3))
	case "local":
		return local.NewStore(cfg.Local)
	default:
		return nil, fmt.Errorf("unrecognized alertmanager storage backend %v, choose one of: azure, configdb, gcs, local, s3", cfg.Type)
	}
}

func newLegacyObjAlertStore(client chunk.ObjectClient, err error) (AlertStore, error) {
	if err != nil {
		return nil, err
	}
	return objectclient.NewAlertStore(client), nil
}
