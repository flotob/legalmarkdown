package main

import (
	"fmt"
	"github.com/eris-ltd/legalmarkdown/lmd"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const CLR_0 = "\x1b[30;1m"
const CLR_R = "\x1b[31;1m"
const CLR_G = "\x1b[32;1m"
const CLR_B = "\x1b[34;1m"
const CLR_N = "\x1b[0m"

func TestLegalToMarkdown(t *testing.T) {
	// create the path properly to the glob command
	test_files_path := filepath.Join(".", "spec", "*.lmd")

	// glob the files
	testfiles, read_error := filepath.Glob(test_files_path)
	if read_error != nil {
		log.Fatal(read_error)
	}

	// set up passed and failed slices
	passed := []string{}
	failed := []string{}

	// run the unit tests
	for _, file := range testfiles {
		success_or_fail := testIndividualFile(file)
		if success_or_fail {
			passed = append(passed, file)
		} else {
			failed = append(failed, file)
			reportResults(passed, failed)
			log.Fatal("Fast fail.")
		}
	}

	reportResults(passed, failed)

}

func testIndividualFile(file string) bool {
	// announce thyself
	fmt.Println(CLR_0, "Testing file: ", file, CLR_N)

	// set the basis and read it into memory
	basis_file := strings.Replace(file, ".lmd", ".md", 1)
	test_against_me := lmd.ReadAFile(basis_file)

	// make a temp file
	temp_file, temp_file_err := ioutil.TempFile(os.TempDir(), "lmd-test-")
	if temp_file_err != nil {
		log.Fatal(temp_file_err)
	}
	defer os.Remove(temp_file.Name())

	// run LegalToMarkdown on the fixture
	LegalToMarkdown(file, "", temp_file.Name())

	// read the tempfile
	i_made_this_file := lmd.ReadAFile(temp_file.Name())

	// announce
	if test_against_me == i_made_this_file {
		fmt.Println(CLR_G, "YES!\n", CLR_N)
		return true
	} else {
		fmt.Println(CLR_R, "NOOOOOOOOOOOOOOOOO.\n", CLR_N)
		fmt.Println(CLR_G, "Expected =>", CLR_N)
		fmt.Println(test_against_me)
		fmt.Println(CLR_R, "Result =>", CLR_N)
		fmt.Println(i_made_this_file)
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
