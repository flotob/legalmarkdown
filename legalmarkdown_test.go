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

func TestLegalToMarkdownWithYAML(t *testing.T) {
	fmt.Println(CLR_B, "\n\tTesting YAML Based Parsing\n", CLR_N)

	// create the path properly to the glob command
	testFilesPath := filepath.Join(".", "spec", "*.lmd")

	// glob the files
	testfiles, readError := filepath.Glob(testFilesPath)
	if readError != nil {
		log.Fatal(readError)
	}

	// set up passed and failed slices
	passed := []string{}
	failed := []string{}

	// run the unit tests
	for _, file := range testfiles {
		successOrFail := testIndividualFileYAML(file)
		if successOrFail {
			passed = append(passed, file)
		} else {
			failed = append(failed, file)
			reportResults(passed, failed)
			log.Fatal("Fast fail.")
		}
	}

	reportResults(passed, failed)

}

func testIndividualFileYAML(file string) bool {
	// announce thyself
	fmt.Println(CLR_0, "Testing file: ", file, CLR_N)

	// set the basis and read it into memory
	basisFile := strings.Replace(file, ".lmd", ".md", 1)
	testAgainstMe := lmd.ReadAFile(basisFile)

	// make a temp file
	tempFile, tempFileErr := ioutil.TempFile(os.TempDir(), "lmd-test-")
	if tempFileErr != nil {
		log.Fatal(tempFileErr)
	}
	defer os.Remove(tempFile.Name())

	// run LegalToMarkdown on the fixture
	lmd.LegalToMarkdown(file, "", tempFile.Name())

	// read the tempfile
	iMadeThisFile := lmd.ReadAFile(tempFile.Name())

	// announce
	if testAgainstMe == iMadeThisFile {
		fmt.Println(CLR_G, "YES!\n", CLR_N)
		return true
	} else {
		fmt.Println(CLR_R, "NOOOOOOOOOOOOOOOOO.\n", CLR_N)
		fmt.Println(CLR_G, "Expected =>", CLR_N)
		fmt.Println(testAgainstMe)
		fmt.Println(CLR_R, "Result =>", CLR_N)
		fmt.Println(iMadeThisFile)
		return false
	}

}

func TestLegalToMarkdownWithJSON(t *testing.T) {
	fmt.Println(CLR_B, "\n\tTesting JSON Based Parsing\n", CLR_N)

	// create the path properly to the glob command
	testFilesPath := filepath.Join(".", "spec", "json", "*.lmd")

	// glob the files
	testfiles, readError := filepath.Glob(testFilesPath)
	if readError != nil {
		log.Fatal(readError)
	}

	// set up passed and failed slices
	passed := []string{}
	failed := []string{}

	// run the unit tests
	for _, file := range testfiles {
		successOrFail := testIndividualFileJSON(file)
		if successOrFail {
			passed = append(passed, file)
		} else {
			failed = append(failed, file)
			reportResults(passed, failed)
			log.Fatal("Fast fail.")
		}
	}

	reportResults(passed, failed)

}

func testIndividualFileJSON(file string) bool {
	// announce thyself
	fmt.Println(CLR_0, "Testing file: ", file, CLR_N)

	// set the basis and read it into memory
	basisFile := strings.Replace(file, ".lmd", ".md", 1)
	paramsFile := strings.Replace(file, ".lmd", ".json", 1)
	testAgainstMe := lmd.ReadAFile(basisFile)

	// make a temp file
	tempFile, tempFileErr := ioutil.TempFile(os.TempDir(), "lmd-test-")
	if tempFileErr != nil {
		log.Fatal(tempFileErr)
	}
	defer os.Remove(tempFile.Name())

	// run LegalToMarkdown on the fixture
	lmd.LegalToMarkdown(file, paramsFile, tempFile.Name())

	// read the tempfile
	iMadeThisFile := lmd.ReadAFile(tempFile.Name())

	// announce
	if testAgainstMe == iMadeThisFile {
		fmt.Println(CLR_G, "YES!\n", CLR_N)
		return true
	} else {
		fmt.Println(CLR_R, "NOOOOOOOOOOOOOOOOO.\n", CLR_N)
		fmt.Println(CLR_G, "Expected =>", CLR_N)
		fmt.Println(testAgainstMe)
		fmt.Println(CLR_R, "Result =>", CLR_N)
		fmt.Println(iMadeThisFile)
		return false
	}

}

func TestLegalToMarkdownHeaders(t *testing.T) {
	fmt.Println(CLR_B, "\n\tTesting Make YAML Front Matter\n", CLR_N)

	// create the path properly to the glob command
	testFilesPath := filepath.Join(".", "spec", "*.lmd")

	// glob the files
	testfiles, readError := filepath.Glob(testFilesPath)
	if readError != nil {
		log.Fatal(readError)
	}

	// set up passed and failed slices
	passed := []string{}
	failed := []string{}

	// run the unit tests
	for _, file := range testfiles {
		successOrFail := testIndividualFileHeaders(file)
		if successOrFail {
			passed = append(passed, file)
		} else {
			failed = append(failed, file)
			reportResults(passed, failed)
			log.Fatal("Fast fail.")
		}
	}

	reportResults(passed, failed)

}

func testIndividualFileHeaders(file string) bool {
	// announce thyself
	fmt.Println(CLR_0, "Testing file: ", file, CLR_N)

	// set the basis and read it into memory
	basisFile := strings.Replace(file, ".lmd", ".headers", 1)
	testAgainstMe := lmd.ReadAFile(basisFile)

	// make a temp file
	tempFile, tempFileErr := ioutil.TempFile(os.TempDir(), "lmd-test-")
	if tempFileErr != nil {
		log.Fatal(tempFileErr)
	}
	defer os.Remove(tempFile.Name())

	// run LegalToMarkdown on the fixture
	lmd.MakeYAMLFrontMatter(file, "", tempFile.Name())

	// read the tempfile
	iMadeThisFile := lmd.ReadAFile(tempFile.Name())

	// announce
	if testAgainstMe == iMadeThisFile {
		fmt.Println(CLR_G, "YES!\n", CLR_N)
		return true
	} else {
		fmt.Println(CLR_R, "NOOOOOOOOOOOOOOOOO.\n", CLR_N)
		fmt.Println(CLR_G, "Expected =>", CLR_N)
		fmt.Println(testAgainstMe)
		fmt.Println(CLR_R, "Result =>", CLR_N)
		fmt.Println(iMadeThisFile)
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
