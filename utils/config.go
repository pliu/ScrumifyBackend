package utils

import (
	"encoding/json"
	"io/ioutil"
)

const path = "./config.json"

type Config struct {
	DB_NAME        string `json:"DB_NAME"`
	DB_PASSWORD    string `json:"DB_PASSWORD"`
	ADMIN_PASSWORD string `json:"ADMIN_PASSWORD"`
}

var Conf Config

func init() {
	data, _ := ioutil.ReadFile(path)
	_ = json.Unmarshal(data, &Conf)
}
