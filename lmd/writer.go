package lmd

import (
    "os"
    "log"
    "io/ioutil"
)

// WriteAFile is a convenience function for writing files.
func WriteAFile(file_to_write string, contents_to_write string) bool {

    // convert to byte array for writing
    contents_as_byte_array := []byte(contents_to_write)

    // if "-" is passed as the file to write, wtite to from stdout instead.
    if file_to_write == " -" || file_to_write == "-" {
        os.Stdout.Write(contents_as_byte_array)
        return true
    }

    // write whole the body
    file_write_err := ioutil.WriteFile(file_to_write, contents_as_byte_array, 0644)
    if file_write_err != nil {
        log.Fatal(file_write_err)
    }
    return true
}