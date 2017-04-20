# mssqlodbc
A simple helper library for GO MSSQL ODBC connections.  This has been tested with 
https://github.com/alexbrainman/odbc


```go
cxn := Connection{
    Server:              "localhost\\SQL2014",
    Database:            "tempdb",
    AppName:             "gosql",
    Trusted:             true,
    MultiSubnetFailover: true,
}

s, err := cxn.ConnectionString()
if err != nil {
    return err
}

db, err := sql.Open("odbc", s)
if err != nil {
    return err
}
defer db.Close()
```

Includes the following features:
* List the valid drivers installed
* Select the "best" driver based on my subjective ranking of them
* Parse a SQL Server ODBC connection string into a Connection

Includes support for: 
* The generic SQL Server ODBC driver
* The SQL Server ODBC Driver v11 and v13
* The SQL Server Native Client ODBC driver v10 and v11

