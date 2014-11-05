package lmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

// writeAFile is a convenience function for writing files. It also does the final cleanup
// by cleaning extraneous new lines and after that parsing and inserting a signature block
// if such is requested by the user.
func writeAFile(file_to_write string, contents_to_write string) bool {

	// close up extraneous new lines
	contents_to_write = strings.Replace(contents_to_write, "\n\n\n", "\n\n", -1)

	signatureBlock := regexp.MustCompile(`\@signature\((.*?):(.*?)\)`)
	if signatureBlock.MatchString(contents_to_write) {

		party1 := signatureBlock.FindAllStringSubmatch(contents_to_write, -1)[0][1]
		party2 := signatureBlock.FindAllStringSubmatch(contents_to_write, -1)[0][2]

		daBlock := fmt.Sprintf(`


______________________________________
Signed: %v



______________________________________
Date



______________________________________
Signed: %v



______________________________________
Date
`, party1, party2)

		contents_to_write = signatureBlock.ReplaceAllString(contents_to_write, daBlock)
	}

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
