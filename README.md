# README

* Command-line flags

        -config=<path>          default: ./config.json
        -env=<test/dev/prod>    Overrides ENV in the config file
        -port=<port>            Overrides PORT in the config file
        
* Configuration file fields

        ENV            default: dev
        PORT           default: 8080 (between 1024 and 49151, inclusive)
        DB_USERNAME    required
        DB_PASSWORD    required
        ADMIN_USERNAME
        ADMIN_PASSWORD
        
* Environments

        test: recreates todo_test on every run and displays both database queries and web requests
        dev: persists todo_dev between runs and only displays web requests
        prod: persists todo_prod between runs and displays neither

* Build

        Windows: go build
        Linux: GOOS=linux GOARCH=amd64 go build