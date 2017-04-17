# mssqlodbc
A simple helper library for GO MSSQL ODBC connections

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
    t.Error("cxn string", err)
}

db, err := sql.Open("odbc", s)
if err != nil {
    t.Error("open", err)
}
defer db.Close()
```

Includes the following features
* Lists the valid drivers installed
* Selects the "best" driver based on my subjective ranking of them

Includes support for 

* The generic SQL Server ODBC driver
* The SQL Server ODBC Driver v11 and v13
* The SQL Server Native Client ODBC driver v10 and v11

