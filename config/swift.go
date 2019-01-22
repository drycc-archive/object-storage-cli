package config

import (
	"github.com/docker/distribution/registry/storage/driver"
	"github.com/docker/distribution/registry/storage/driver/factory"
	// this blank import is used to register the swift driver with the storage driver factory
	_ "github.com/docker/distribution/registry/storage/driver/swift"
)

// Swift is the Config implementation for the Swift client
type Swift struct {
	UserName      string `envconfig:"USER_NAME_FILE" default:"/var/run/secrets/drycc/objectstore/creds/username"`
	Password      string `envconfig:"PASSWORD_FILE" default:"/var/run/secrets/drycc/objectstore/creds/password"`
	ContainerFile string `envconfig:"CONTAINER_FILE" default:"/var/run/secrets/drycc/objectstore/creds/container"`
	AuthURL       string `envconfig:"AUTHURL_FILE" default:"/var/run/secrets/drycc/objectstore/creds/authurl"`
	Tenant        string `envconfig:"TENANT_FILE" default:"/var/run/secrets/drycc/objectstore/creds/tenant"`
	AuthVersion   string `envconfig:"AUTH_VERSION_FILE" default:"/var/run/secrets/drycc/objectstore/creds/authversion"`
}

// CreateDriver is the Config interface implementation
func (s Swift) CreateDriver() (driver.StorageDriver, error) {
	files, err := readFiles(true, s.UserName, s.Password, s.ContainerFile, s.AuthURL)
	if err != nil {
		return nil, err
	}
	username, password, container, authurl := files[0], files[1], files[2], files[3]
	params := map[string]interface{}{
		"username":  username,
		"password":  password,
		"container": container,
		"authurl":   authurl,
	}
	//tenant and authversion are not mandatory and hence need not panc if not present
	tenantfile, err := readFiles(true, s.Tenant)
	if err == nil {
		params["tenant"] = tenantfile[0]
	}
	authversionfile, err := readFiles(true, s.AuthVersion)
	if err == nil {
		params["authversion"] = authversionfile[0]
	}
	return factory.Create("swift", params)
}

// Name is the fmt.Stringer interface implementation
func (s Swift) String() string {
	return SwiftStorageType.String()
}
