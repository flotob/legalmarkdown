package lmd

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"log"
	"regexp"
	"strconv"
)

// HandleParameterAssembly is the primary parsing function which is used by the cli
// assemble method and any libraries seek to build YAML front matter for a template
// file.
//
// The function accepts a string of content and a map of parameters.
//
// First the function will establish four maps, one each for the mixins, the optional
// clauses, the headers, and they style parameters. Then the function will run through
// searching functions to handle the assembly of each of these four major elements.
//
// If any parameters have been sent to the function via the paramaters map, the values
// in each of the parameters field will be maintained. If the keys in the parameters
// map do not exist, but the appropriate mixin, optional clause, or header exists in the
// text of the template, then the appropriate key will be added to the appropriate map.
//
// Lastly the function will call a function that will reassemble the file by running
// through each of the maps and building the front matter. The reassembled contents
// are returned to the calling function.
func HandleParameterAssembly(contents string, parameters map[string]string) string {

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

// HandleParameterAssemblyJSON performs roughly the same parsing function as the
// HandleParamaterAssembly function.
//
// The major difference is that once the four maps are established by the search functions
// then each of the keys and values in these maps are copied back into the parameters
// map.
//
// Finally the reassembled parameters map is marshaled into a JSON string that is
// returned to the calling function.
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

// findTheMixins runs through the text of the content to find mixin patterns.
// These are then placed into the mixins map, which is rationalized against the
// parameters map so that the values which are passed from the parameters map
// when the contents are originally read into memory are not overwritten.
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

// findTheOptClauses performs exactly the same function as the findTheMixins function
// except it is parsing the text for the optional clauses pattern instead of the mixins
// pattern.
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

// findTheLeaders runs through the text first to determine if there is a block. If there is no
// block then it returns the contents and empty maps to the calling function. If there is a block
// then the function uses the splitTheBlock function to gain a string of all of the headers.
//
// After checking whether the header style is oldStyle ("llll.") or newStyle ("l4.") then the
// function loops through the slice and for each of the leaders it sinks these into the headers
// map which also checking if the value from the parameters map is kept.
//
// Finally the function assembles a three length map for the styles by performing roughly the
// same algorithm as the rest of this file to ensure that the values of the parameters map are maintained.
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

// assembleStyle is a convenience function which simply checks if the passed style
// parameter is already in the parameters map and if so, it sinks that value into the
// style map and returns the map.
func assembleStyle(style string, parameters map[string]string, styles map[string]string) map[string]string {
	if _, exists := parameters[style]; !exists {
		styles[style] = ""
	} else {
		styles[style] = parameters[style]
	}
	return styles
}

// reAssembleTheFile is a convenience function which builds the front matter in a particular way
// which will make it easy for users of the legalmarkdown system to understand how to use the system.
//
// It first will build the mixins, then the Optional Clauses, then the headers, and finally the styles
// along the way if any of these blocks does not exist then the appropriate section will not be built.
//
// Before any of that building, the function will check to make sure whether all of the main maps are
// empty (in which case there is no front matter to build and the content without front matter is
// returned to the calling function).
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
