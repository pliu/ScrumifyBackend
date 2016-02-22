package utils

import(
"encoding/json"
"io/ioutil"
)

const path = "./config.json"

type Config struct {
	DB_NAME      string  `json:"DB_NAME"`
	DB_PASSWORD  string  `json:"DB_PASSWORD"`
}

var config Config

func init() {
	data, _ := ioutil.ReadFile(path)
	_ = json.Unmarshal(data, &config)
}

func GetConfig() Config {
	return config;
}