package main

import (
	"encoding/json"
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
		t.Error(readError)
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
			t.Error("Fast fail.")
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
		t.Error(readError)
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
			t.Error("Fast fail.")
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
		t.Error(readError)
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
			t.Error("Fast fail.")
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

func TestGetParameters(t *testing.T) {
	fmt.Println(CLR_B, "\n\tTesting Get Parameters\n", CLR_N)

	basisFile := filepath.Join(".", "spec", "json", "30.block_all_leader_types.json")
	param := make(map[string]string)
	file_buffer, _ := ioutil.ReadFile(basisFile)
	json.Unmarshal(file_buffer, &param)
	paramsAsJsonByteArray, _ := json.Marshal(param)
	basis := string(paramsAsJsonByteArray)

	testFile1 := filepath.Join(".", "spec", "30.block_all_leader_types.lmd")
	test1 := lmd.GetTheParameters(testFile1)

	if test1 != basis {
		fmt.Println(CLR_R, "NOOOOOOOOOOOOOOOOO.\n", CLR_N)
		fmt.Println(CLR_G, "Expected =>", CLR_N)
		fmt.Println(basis)
		fmt.Println(CLR_R, "Result =>", CLR_N)
		fmt.Println(test1)
		t.Error("JSONizing parameters #1 -- from YAML front matter -- failed.\n")
	} else {
		fmt.Println(CLR_G, "JSONizing parameters #1 -- from YAML front matter -- passed.\n", CLR_N)
	}

	testFile2 := filepath.Join(".", "spec", "json", "30.block_all_leader_types.lmd")
	test2 := lmd.GetTheParameters(testFile2)
	test2 = strings.Replace(test2, `level-style":"",`, "", 1)
	test2 = strings.Replace(test2, `,"no-reset":""`, "", 1)
	test2 = strings.Replace(test2, `""no-indent":""`, `"no-indent":""`, 1)

	for k, _ := range param {
		param[k] = ""
	}
	paramsAsJsonByteArray, _ = json.Marshal(param)
	basis = string(paramsAsJsonByteArray)

	if test2 != basis {
		fmt.Println(CLR_R, "NOOOOOOOOOOOOOOOOO.\n", CLR_N)
		fmt.Println(CLR_G, "Expected =>", CLR_N)
		fmt.Println(basis)
		fmt.Println(CLR_R, "Result =>", CLR_N)
		fmt.Println(test2)
		t.Error("JSONizing parameters #2 -- made from a raw LMD -- failed.\n")
	} else {
		fmt.Println(CLR_G, "JSONizing parameters #2 -- made from a raw LMD -- passed.\n", CLR_N)
	}
}

func TestLegalToRenderingToPDF(t *testing.T) {
	fmt.Println(CLR_B, "\n\tTesting Rendering to PDF\n", CLR_N)

	testFile := filepath.Join(".", "spec", "00.load_write_no_action.md")

	// make a temp file
	tempFile, tempFileErr := ioutil.TempFile(os.TempDir(), "lmd-test-")
	if tempFileErr != nil {
		t.Error(tempFileErr)
	}
	defer os.Remove(tempFile.Name())

	// dowit
	lmd.MarkdownToPDF(testFile, "", tempFile.Name())

	// read the tempfile
	iMadeThisFile := lmd.ReadAFile(tempFile.Name())

	if iMadeThisFile == "" {
		t.Error("Did not create a pdf.")
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
