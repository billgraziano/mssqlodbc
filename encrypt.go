package mssqlodbc

import "github.com/pkg/errors"

// https://learn.microsoft.com/en-us/troubleshoot/sql/connect/certificate-chain-not-trusted?tabs=odbc-driver-18x
// https://techcommunity.microsoft.com/t5/sql-server-blog/odbc-driver-18-0-for-sql-server-released/ba-p/3169228
// https://learn.microsoft.com/en-us/sql/connect/odbc/windows/release-notes-odbc-sql-server-windows?view=sql-server-ver16

var ErrInvalidEncrypt = errors.New("invalid encrypt: expected (blank), Optional, Yes, No, Mandatory, Strict")

const (
	EncryptYes       string = "Yes"
	EncryptMandatory string = "Mandatory"
	EncryptNo        string = "No"
	EncryptOptional  string = "Optional"
	EncryptStrict    string = "Strict"
)
