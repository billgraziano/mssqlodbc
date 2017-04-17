package mssqlodbc

import (
	"fmt"

	"github.com/pkg/errors"
)

// Connection holds information about an ODBC SQL Server connection
type Connection struct {
	driver              string
	Server              string
	User                string
	Password            string
	Trusted             bool
	AppName             string
	Database            string
	MultiSubnetFailover bool
}

// Driver gets the driver for the connection
func (c *Connection) Driver() string {
	return c.driver
}

// SetDriver sets the driver for a connection
func (c *Connection) SetDriver(d string) error {
	err := ValidDriver(d)
	if err != nil {
		return err
	}
	c.driver = d
	return nil
}

// ConnectionString returns a connection string
func (c *Connection) ConnectionString() (string, error) {

	// https://docs.microsoft.com/en-us/sql/relational-databases/native-client/applications/using-connection-string-keywords-with-sql-server-native-client

	var cxn string

	if c.driver == "" {
		driver, err := BestDriver()
		if err == ErrNoDrivers {
			return "", err
		}
		if err != nil {
			return "", errors.Wrap(err, "bestdriver")
		}
		c.driver = driver
	}

	// Driver
	cxn += fmt.Sprintf("Driver={%s}; ", c.driver)

	// Host
	// {SQL Server needs Server so we'll use this as default}
	if c.Server == "" {
		return "", errors.New("invalid server")
	}
	cxn += fmt.Sprintf("Server=%s; ", c.Server)

	// Authentication
	if c.Trusted || (c.User == "" && c.Password == "") {
		cxn += fmt.Sprintf("Trusted_Connection=yes; ")
	} else {
		cxn += fmt.Sprintf("UID=%s; PWD=%s; ", c.User, c.Password)
	}

	// Database
	if c.Database != "" {
		cxn += fmt.Sprintf("Database=%s; ", c.Database)
	}

	// App Name
	if c.AppName != "" {
		cxn += fmt.Sprintf("App=%s; ", c.AppName)
	}

	// MultisubnetFailover
	if c.MultiSubnetFailover {
		cxn += "MultiSubnetFailover=Yes; "
	}

	return cxn, nil
}
