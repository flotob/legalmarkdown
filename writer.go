package legalmarkdown

import (
    "log"
    "io/ioutil"
)

func write_a_file(file_to_write string, contents_to_write string) bool {
    // TODO: need stdin gaurd if file_to_write == '-'

    contents_as_byte_array := []byte(contents_to_write)

    // write whole the body
    file_write_err := ioutil.WriteFile(file_to_write, contents_as_byte_array, 0644)
    if file_write_err != nil {
        log.Fatal(file_write_err)
    }

    return true
}