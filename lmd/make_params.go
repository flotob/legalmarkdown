package lmd

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"log"
	"regexp"
	"strconv"
)

// HandleParameterAssembly ...
func HandleParameterAssembly(contents string, parameters map[string]string) string {

	if len(parameters) == 0 {
		return contents
	}

	mixins := make(map[string]string)
	optClauses := make(map[string]string)
	headers := make(map[string]string)
	styles := make(map[string]string)

	contents, mixins = findTheMixins(contents, parameters)
	contents, optClauses = findTheOptClauses(contents, parameters)
	contents, headers, styles = findTheLeaders(contents, parameters)
	contents = reAssembleTheFile(contents, mixins, optClauses, headers, styles)

	return contents

}

// HandleParameterAssemblyJSON ...
func AssembleParametersIntoJSON(contents string, parameters map[string]string) string {

	mixins := make(map[string]string)
	optClauses := make(map[string]string)
	headers := make(map[string]string)
	styles := make(map[string]string)

	contents, mixins = findTheMixins(contents, parameters)
	contents, optClauses = findTheOptClauses(contents, parameters)
	contents, headers, styles = findTheLeaders(contents, parameters)

	for k, v := range mixins {
		parameters[k] = v
	}
	for k, v := range optClauses {
		parameters[k] = v
	}
	for k, v := range headers {
		parameters[k] = v
	}
	for k, v := range styles {
		parameters[k] = v
	}

	paramsAsJsonByteArray, err := json.Marshal(parameters)

	if err != nil {
		log.Fatal("JSON assembly error.")
	}

	return string(paramsAsJsonByteArray)

}

// findTheMixins ...
func findTheMixins(contents string, parameters map[string]string) (string, map[string]string) {

	mixins := make(map[string]string)
	mixinPattern := regexp.MustCompile(`[^\[]{{(\S+)}}`)

	if mixinPattern.MatchString(contents) {
		for _, matchSlice := range mixinPattern.FindAllStringSubmatch(contents, -1) {
			if _, exists := parameters[matchSlice[1]]; !exists {
				mixins[matchSlice[1]] = ""
			} else {
				mixins[matchSlice[1]] = parameters[matchSlice[1]]
			}
		}
	}

	return contents, mixins
}

// findTheOptClauses ...
func findTheOptClauses(contents string, parameters map[string]string) (string, map[string]string) {

	optClauses := make(map[string]string)
	optClausesPattern := regexp.MustCompile(`\[{{(\S+)}}`)

	if optClausesPattern.MatchString(contents) {
		for _, matchSlice := range optClausesPattern.FindAllStringSubmatch(contents, -1) {
			if _, exists := parameters[matchSlice[1]]; !exists {
				optClauses[matchSlice[1]] = ""
			} else {
				optClauses[matchSlice[1]] = parameters[matchSlice[1]]
			}
		}
	}

	return contents, optClauses
}

// findTheLeaders ...
func findTheLeaders(contents string, parameters map[string]string) (string, map[string]string, map[string]string) {

	headers := make(map[string]string)
	styles := make(map[string]string)

	to_run, _, block, _ := findTheBlock(contents)
	if !to_run {
		return contents, headers, styles
	}
	_, leadersSlice := splitTheBlock(block)

	oldStyle := true
	for _, l := range leadersSlice {
		if l == "l1." {
			oldStyle = false
		}
	}

	for _, leader := range leadersSlice {
		if oldStyle {
			leader = strconv.Itoa(len(leader) - 1)
			leader = "level-" + leader
		} else {
			leader = string(leader[1])
			leader = "level-" + leader
		}
		if _, exists := parameters[leader]; !exists {
			headers[leader] = ""
		} else {
			headers[leader] = parameters[leader]
		}
	}

	stylesSlice := []string{"no-indent", "no-reset", "level-style"}

	for _, style := range stylesSlice {
		styles = assembleStyle(style, parameters, styles)
	}

	return contents, headers, styles
}

// assembleStyle ...
func assembleStyle(style string, parameters map[string]string, styles map[string]string) map[string]string {
	if _, exists := parameters[style]; !exists {
		styles[style] = ""
	} else {
		styles[style] = parameters[style]
	}
	return styles
}

// reAssembleTheFile ...
func reAssembleTheFile(contents string, mixins map[string]string, optClauses map[string]string, headers map[string]string, styles map[string]string) string {

	if !(len(mixins) == 0) || !(len(optClauses) == 0) || !(len(headers) == 0) {
		frontMatter := "---\n\n"
		if !(len(mixins) == 0) {
			mixinsAsByteArray, _ := yaml.Marshal(mixins)
			mixinString := string(mixinsAsByteArray)
			frontMatter = frontMatter + "# Mixins\n" + mixinString
		}
		if !(len(mixins) == 0) && !(len(optClauses) == 0) {
			frontMatter = frontMatter + "\n"
		}
		if !(len(optClauses) == 0) {
			optClausesAsByteArray, _ := yaml.Marshal(optClauses)
			optClausesString := string(optClausesAsByteArray)
			frontMatter = frontMatter + "# Optional Clauses\n" + optClausesString
		}
		if (!(len(mixins) == 0) || !(len(optClauses) == 0)) && !(len(headers) == 0) {
			frontMatter = frontMatter + "\n"
		}
		if !(len(headers) == 0) {
			headersAsByteArray, _ := yaml.Marshal(headers)
			headersString := string(headersAsByteArray)
			stylesAsByteArray, _ := yaml.Marshal(styles)
			stylesString := string(stylesAsByteArray)
			frontMatter = frontMatter + "# Structured Headers\n" + headersString
			frontMatter = frontMatter + "\n# Properties\n" + stylesString
		}
		frontMatter = frontMatter + "\n---\n\n"
		contents = frontMatter + contents
	}

	return contents
}
