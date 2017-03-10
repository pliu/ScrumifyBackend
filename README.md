# README

* Dependencies

        github.com/gin-gonic/gin (customized; submitted a PR)
        github.com/go-sql-driver/mysql
        github.com/stretchr/testify
        gopkg.in/gorp.v2 (customized)
        
* Customizations (gorp.v2)

        Added to SqlForCreate in table.go:
        if col.DefaultStatement != "" {
        	s.WriteString(" " + col.DefaultStatement)
        }
        
        if len(t.foreignKeys) > 0 {
        	for _, foreignKey := range t.foreignKeys {
        		t := bytes.Buffer{}
        		s.WriteString(", foreign key (")
       			t.WriteString("(")
        		for i, mapping := range foreignKey.mappings {
        			if i > 0 {
        				s.WriteString(", ")
       					t.WriteString(", ")
        			}
        			s.WriteString(mapping.FieldName)
        			t.WriteString(mapping.ReferenceFieldName)
        		}
       			s.WriteString(") references ")
        		t.WriteString(")")
        		s.WriteString(foreignKey.referenceTable)
        		s.WriteString(t.String())
       			if foreignKey.statements != "" {
        			s.WriteString(" ")
        			s.WriteString(foreignKey.statements)
        		}
        	}
       	}
       	
       	Added new method to table.go:
       	func (t *TableMap) SetForeignKeys(referenceTable string, statements string, mappings ...FieldNameMapping) *TableMap {
        	var foreignKey ForeignKey
        	foreignKey.referenceTable = referenceTable
        	foreignKey.statements = statements
        	for _, mapping := range mappings {
        		foreignKey.mappings = append(foreignKey.mappings, mapping)
        	}
        	t.foreignKeys = append(t.foreignKeys, foreignKey)
        
        	return t
        }
        
        Added new file, foreign_key_structs.go:
        type FieldNameMapping struct {
        	FieldName string
        	ReferenceFieldName string
        }
        
        type ForeignKey struct {
        	mappings []FieldNameMapping
        	referenceTable string
        	statements string
        }
        
        Added to ColumnMap in column.go:
        DefaultStatement string
        
        Added new method to column.go:
        func (c *ColumnMap) SetDefaultStatement(str string) *ColumnMap {
        	c.DefaultStatement = str
        	return c
        }
        
        Changed bindInsert in table_bindings.go:
        if !col.Transient -> if !col.Transient && col.DefaultStatement == ""
        
        Changed bindUpdate in table_bindings.go:
        if !col.isAutoIncr && !col.Transient && colFilter(col) -> if !col.isAutoIncr && !col.Transient && colFilter(col) && col.DefaultStatement == ""
        
        Changed ToSqlType in dialect_mysql.go:
        case "Time": -> case "Time", "NullTime":
        
        Added to ToSqlType in dialect_mysql.go:
        case "Dependencies":
        	return "blob"
        	
* Customizations (gin)

        Added to LoggerWithWriter in logger.go:
        if c, err := os.Open("CONOUT$"); err == nil && isatty.IsTerminal(c.Fd()) {
        	isTerm = true
        }

* Command-line flags

        -config=<path>          default: ./config.json
        -cert                   overrides CERT_PATH in the config file
        -key                    overrides KEY_PATH in the config file
        -env=<test/dev/prod>    overrides ENV in the config file
        -port=<port>            overrides PORT in the config file
        
* Configuration file fields

        CERT_PATH           default: ./cert.pem
        KEY_PATH            default: ./key.pem
        ENV                 default: dev
        PORT                default: 8080 (between 1024 and 49151, inclusive)
        DB_USERNAME         required
        DB_PASSWORD         required
        ADMIN_USERNAME
        ADMIN_PASSWORD
        
* Environments

        test:   recreates todo_test on every run and only displays web requests
        dev:    persists todo_dev between runs and displays both database queries and web requests
        prod:   persists todo_prod between runs and only displays web requests

* Build

        Windows:    go build
        Linux:      GOOS=linux GOARCH=amd64 go build
        