package mssqlodbc

import (
	"fmt"
	"testing"
)

func TestBestDriver(t *testing.T) {
	s, err := BestDriver()
	if err != nil {
		t.Error("Best Driver error: ", err)
	}
	fmt.Println("Best Driver: ", s)
}
func TestODBCDriver(t *testing.T) {

	d := NativeClient11
	fmt.Println("Native 11: ", d)
}

func TestEmptyDriver(t *testing.T) {
	var d string
	fmt.Println("Empty Driver: ", NoDriver)
	if d != NoDriver {
		t.Error("empty string != NoDriver")
	}
}

func TestAvailableDrivers(t *testing.T) {
	fmt.Println("Available Drivers")
	fmt.Println("=====================================")
	d, err := AvailableDrivers()
	if err != nil {
		t.Error("available drivers: ", err)
	}
	for _, s := range d {
		fmt.Println(s)
	}
}

func TestValidDrivers(t *testing.T) {
	_, err := ValidDriver("test")
	if err != ErrInvalidDriver {
		t.Error("'test' should be invalid driver")
	}

	v, err := ValidDriver("SQL Server Native Client 11.0")
	if v == false || err != nil {
		t.Error("Native11: failed")
	}
}
