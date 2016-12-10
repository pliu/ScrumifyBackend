package utils

import (
    "log"
)

func FatalErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}

func PrintErr(err error, msg string) {
    if err != nil {
        log.Println(msg, err)
    }
}
