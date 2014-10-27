package legalmarkdown

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

// read_a_file is a convenience function. given a filename string, reads the file and passes it back to
// the calling function as a string.
func read_a_file(file_to_read string) string {
	// TODO: need stdin gaurd if file_to_read == '-'

	file_buffer, file_read_err := ioutil.ReadFile(file_to_read)

	if file_read_err != nil {
		log.Fatal(file_read_err)
	}

	contents := string(file_buffer)
	return contents
}

// import_files. included files need to be imported into the primary contents string. this function manages
// that process. first it compiles a regular expression which will search for the trigger string
//   @include PARTIAL
//
// if one or more match is found, the function will simply replace the `@include PARTIAL` line with the
// read in string of the included partial. the complete string will be returned to the calling function.
func import_files(fileContents string) string {
	importRegExp := regexp.MustCompile(`(?m)^@include (.*?)$`)

	if importRegExp.MatchString(fileContents) {
		importedFiles := importRegExp.FindAllStringSubmatch(fileContents, -1)
		for _, importedFile := range importedFiles {
			fileContents = strings.Replace(fileContents, importedFile[0], read_a_file(importedFile[1]), -1)
		}
		return fileContents
	} else {
		return fileContents
	}
}

// parse_template_to_find_paramamters. paramaters may be passed to the parser either separately from the
// template file or as part of the template file. this function manages the process of stripping paramaters
// out of a template file. the function first compiles a YAML Front Matter regular expression. then if a
// match for that regular expression is found, the contents of the template file are replaced with an empty
// string and the YAML front matter is returned, along with the replaced contents (both as strings) to the
// calling function.
func parse_template_to_find_parameters(fileContents string) (string, string) {
	yamlRegExp := regexp.MustCompile(`(?sm)\A(---\s*\n.*?)(^---\s*\n)`)

	if yamlRegExp.MatchString(fileContents) {
		yamlFrontMatter := yamlRegExp.FindAllStringSubmatch(fileContents, -1)[0][1]
		fileContents = yamlRegExp.ReplaceAllString(fileContents, "")
		return yamlFrontMatter, fileContents
	} else {
		return "", fileContents
	}
}

// unmarshall_parameters. yaml and json paramaters must be unmarshalled into the paramaters map. this function
// is responsible for unmarshalling the paramaters from yaml or json strings into (first a byte array) and
// subsequently into the paramaters map which is returned to the calling function.
func unmarshall_parameters(parameters string) map[string]string {
	// TODO: make this smarter... should be able to also parse JSON if the YAML unmarshall fails
	parameter_bytes := []byte(parameters)
	param := make(map[string]string)
	yaml.Unmarshal(parameter_bytes, &param)
	return param
}

// merge_parameters. this is a convenience function which will merge two hash maps. any conflicting parameters
// in the two maps will be resolved in favor of the *first* map which is passed. that is to say that the first
// map passed to the function, the `superior_map` map, will overwrite the `sublimated_map`.
func merge_parameters(superior_map map[string]string, sublimated_map map[string]string) map[string]string {
	for k, v := range superior_map {
		sublimated_map[k] = v
	}
	return sublimated_map
}
