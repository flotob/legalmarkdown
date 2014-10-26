package main

import (
    "log"
    "regexp"
    "strings"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

func read_a_file(file_to_read string) (string) {
    // TODO: need stdin gaurd if file_to_read == '-'
    file_buffer, file_read_err := ioutil.ReadFile(file_to_read)

    if file_read_err != nil {
        log.Fatal(file_read_err)
    }

    contents := string(file_buffer)

    return contents
}

func import_files(fileContents string) string {
    importRegExp := regexp.MustCompile(`(?m)^@include (.*?)$`)

    if importRegExp.MatchString(fileContents) {
        importedFiles       := importRegExp.FindAllStringSubmatch(fileContents, -1)
        for _, importedFile := range importedFiles {
            fileContents = strings.Replace(fileContents, importedFile[0], read_a_file(importedFile[1]), -1)
        }
        return fileContents
    } else {
        return fileContents
    }
}

func parse_template_to_find_parameters(fileContents string) (string, string) {
    yamlRegExp := regexp.MustCompile(`(?sm)\A(---\s*\n.*?)(^---\s*\n)`)

    if yamlRegExp.MatchString(fileContents) {
        yamlFrontMatter  := yamlRegExp.FindAllStringSubmatch(fileContents, -1)[0][1]
        fileContents      = yamlRegExp.ReplaceAllString(fileContents, "")
        return yamlFrontMatter, fileContents
    } else {
        return "", fileContents
    }
}

func unmarshall_parameters(parameters string) (map[string]string) {
    parameter_bytes := []byte(parameters)
    param           := make(map[string]string)
    yaml.Unmarshal([]byte(parameter_bytes), &param)
    return param
}

func merge_parameters(amended_parameters map[string]string, merged_parameters map[string]string) (map[string]string) {
    for k, v := range amended_parameters {
        merged_parameters[k] = v
    }
    return merged_parameters
}