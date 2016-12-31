package utils

import (
    "encoding/json"
    "io/ioutil"
    "flag"
    "errors"
)

type Config struct {
    CERT_PATH      string `json:"CERT_PATH"`
    KEY_PATH       string `json:"KEY_PATH"`
    ENV            string `json:"ENV"`
    PORT           int64  `json:"PORT"`
    DB_USERNAME    string `json:"DB_USERNAME"`
    DB_PASSWORD    string `json:"DB_PASSWORD"`
    ADMIN_USERNAME string `json:"ADMIN_USERNAME"`
    ADMIN_PASSWORD string `json:"ADMIN_PASSWORD"`
}

var Conf Config

func InitializeConfig() {
    var config_path string
    flag.StringVar(&config_path, "config", "./config.json", "Path to the config file; defaults to ./config.json")
    var cert_path string
    flag.StringVar(&cert_path, "cert", "", "Path to the cert file; overrides config file and defaults to ./cert.pem")
    var key_path string
    flag.StringVar(&key_path, "key", "", "Path to the key file; overrides config file and defaults to ./key.pem")
    var env string
    flag.StringVar(&env, "env", "", "test, dev, or prod; overrides config file and defaults to dev")
    var port int64
    flag.Int64Var(&port, "port", 0, "Hosting port; overrides config file and defaults to 8080")
    flag.Parse()

    readConfig(config_path)
    if cert_path != "" {
        Conf.CERT_PATH = cert_path
    }
    if Conf.CERT_PATH == "" {
        Conf.CERT_PATH = "./cert.pem"
    }
    if key_path != "" {
        Conf.KEY_PATH = key_path
    }
    if Conf.KEY_PATH == "" {
        Conf.KEY_PATH = "./key.pem"
    }
    if env == "test" {
        Conf.ENV = "test"
    } else if env == "prod" {
        Conf.ENV = "prod"
    } else if env == "dev" || Conf.ENV == "" {
        Conf.ENV = "dev"
    }
    if port != 0 {
        Conf.PORT = port
    }
    if Conf.PORT == 0 {
        Conf.PORT = 8080
    }

    validateConfig()
}

func readConfig(path string) {
    if data, err := ioutil.ReadFile(path); err == nil {
        if err = json.Unmarshal(data, &Conf); err != nil {
            FatalErr(err, "Parse config failed")
        }
    } else {
        FatalErr(err, "Read config at " + path + " failed")
    }
}

func validateConfig() {
    if Conf.PORT < 1024 || Conf.PORT > 49151 || Conf.DB_USERNAME == "" || Conf.DB_PASSWORD == "" {
        FatalErr(errors.New("Check config"), "Invalid configration")
    }
}
