# README

* Dependencies

        github.com/gin-gonic/gin
        github.com/go-sql-driver/mysql
        github.com/pjebs/restgate
        github.com/unrolled/secure
        gopkg.in/gorp.v2 (customized)
        
* Customizations

        Added to SqlForCreate in table.go:
        if col.DefaultStatement != "" {
        	s.WriteString(" " + col.DefaultStatement)
        }
        
        Added to ColumnMap in column.go:
        DefaultStatement string
        
        Changed bindInsert in table_bindings.go:
        if !col.Transient -> if !col.Transient && col.DefaultStatement == ""
        
        Changed bindUpdate in table_bindings.go:
        if !col.isAutoIncr && !col.Transient && colFilter(col) -> if !col.isAutoIncr && !col.Transient && colFilter(col) && col.DefaultStatement == ""
        
        Changed ToSqlType in dialect_mysql.go:
        case "Time": -> case "Time", "NullTime":

* Command-line flags

        -config=<path>          default: ./config.json
        -env=<test/dev/prod>    overrides ENV in the config file
        -port=<port>            default: 8080; overrides PORT in the config file
        
* Configuration file fields

        ENV                 default: dev
        PORT                default: 8080 (between 1024 and 49151, inclusive)
        DB_USERNAME         required
        DB_PASSWORD         required
        ADMIN_USERNAME
        ADMIN_PASSWORD
        
* Environments

        test:   recreates todo_test on every run and displays both database queries and web requests
        dev:    persists todo_dev between runs and only displays web requests
        prod:   persists todo_prod between runs and displays neither

* Build

        Windows:    go build
        Linux:      GOOS=linux GOARCH=amd64 go build
        