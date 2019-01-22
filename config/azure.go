package config

import (
	"github.com/docker/distribution/registry/storage/driver"
	"github.com/docker/distribution/registry/storage/driver/factory"
	// this blank import is used to register the Azure driver with the storage driver factory
	_ "github.com/docker/distribution/registry/storage/driver/azure"
)

// Azure is the Config implementation for the Azure client
type Azure struct {
	AccountNameFile string `envconfig:"ACCOUNT_NAME_FILE" default:"/var/run/secrets/drycc/objectstore/creds/accountname"`
	AccountKeyFile  string `envconfig:"ACCOUNT_KEY_FILE" default:"/var/run/secrets/drycc/objectstore/creds/accountkey"`
	ContainerFile   string `envconfig:"CONTAINER_FILE" default:"/var/run/secrets/drycc/objectstore/creds/container"`
}

// CreateDriver is the Config interface implementation
func (a Azure) CreateDriver() (driver.StorageDriver, error) {
	files, err := readFiles(true, a.AccountNameFile, a.AccountKeyFile, a.ContainerFile)
	if err != nil {
		return nil, err
	}
	accountName, accountKey, container := files[0], files[1], files[2]
	params := map[string]interface{}{
		"accountname": accountName,
		"accountkey":  accountKey,
		"container":   container,
	}
	return factory.Create("azure", params)
}

// Name is the fmt.Stringer interface implementation
func (a Azure) String() string {
	return AzureStorageType.String()
}
