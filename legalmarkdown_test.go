package legalmarkdown

import (
    "os"
    "fmt"
    "log"
    "strings"
    "testing"
    "io/ioutil"
    "path/filepath"
)

const CLR_0 = "\x1b[30;1m"
const CLR_R = "\x1b[31;1m"
const CLR_G = "\x1b[32;1m"
const CLR_B = "\x1b[34;1m"
const CLR_N = "\x1b[0m"

func TestLegalToMarkdown(t *testing.T) {
    // create the path properly to the glob command
    test_files_path       := filepath.Join(".", "spec", "*.lmd")

    // glob the files
    testfiles, read_error := filepath.Glob(test_files_path)
    if read_error != nil {
        log.Fatal(read_error)
    }

    // set up passed and failed slices
    passed := []string{}
    failed := []string{}

    // run the unit tests
    for _, file := range testfiles  {
        success_or_fail := testIndividualFile(file)
        if success_or_fail {
            passed = append(passed, file)
        } else {
            failed = append(failed, file)
        }
    }

    reportResults(passed, failed)

    if len(failed) != 0 {
        t.Error("Did Not Pass the Tests.")
    }
}

func testIndividualFile(file string) bool {
    // announce thyself
    fmt.Println(CLR_0, "Testing file: ", file, CLR_N)

    // set the basis and read it into memory
    basis_file      := strings.Replace(file, ".lmd", ".md", 1)
    test_against_me := read_a_file(basis_file)

    // make a temp file
    temp_file, temp_file_err := ioutil.TempFile(os.TempDir(), "lmd-test-")
    if temp_file_err != nil {
        log.Fatal(temp_file_err)
    }
    defer os.Remove(temp_file.Name())

    // run LegalToMarkdown on the fixture
    LegalToMarkdown(file, "", temp_file.Name())

    // read the tempfile
    i_made_this_file := read_a_file(temp_file.Name())

    // announce
    if test_against_me == i_made_this_file {
        fmt.Println(CLR_G, "YES!\n", CLR_N)
        return true
    } else {
        fmt.Println(CLR_R, "NOOOOOOOOOOOOOOOOO.\n", CLR_N)
        return false
    }

}

func reportResults(passed []string, failed []string) {

    fmt.Println("")
    fmt.Println(CLR_B, "*****", CLR_N)
    fmt.Println("")
    fmt.Println(CLR_G, "Tests Passed: ", len(passed), CLR_N)
    fmt.Println(CLR_R, "Tests Failed: ", len(failed), CLR_N)
    fmt.Println("")
    fmt.Println(CLR_B, "*****", CLR_N)
    fmt.Println("")

}