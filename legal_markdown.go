package main

import (
    // "os"
    "fmt"
    "log"
    "io/ioutil"
    // "path/filepath"
    // "strings"
)

func main(){

    buff, err := ioutil.ReadFile("spec/00.load_write_no_action.lmd")

    if err != nil {
        log.Fatal(err)
    }

    contents := string(buff)
    fmt.Print(contents)

}