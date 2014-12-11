package lmd

// LegalToMarkdown is the primary function which controls parsing a template document into a markdown
// result when the parsing library is called from the command line. Two strings which are filenames
// should be passed to the function. The parameters string may be an empty string. The function first
// parses and reads the command sent from the command line, and then reads the template file. After
// this, the function pulls into the template file any partials which have been included into the
// template with the `@include {{PARTIAL}}` flag within the text of the primary template file that
// has been called.
//
// Then the function reads the paramaters from a parameters file, the template file, or both. In the
// case where the parameters are contained in both a parameters file and in the template file, the
// parameters in the template file are considered as defaults which are overridden by parameters passed
// to the function via the paramaters file.
//
// Finally once the function has prepared the `contents` and `parameters` from the various passed files
// and built a cohesive set of `contents` and `parameters`.
//
// These are passed to the primary entrance function to the parsing process. The contents of the file
// which needs to be parsed, and the parameters which should control the parsing and transformation
// of the lmd file to a rendered document are lexed and ready for the parser to run through
// the sequence of mixins, optional clauses, and structured headers.
//
// The parser will first call the primary mixins function, then will call the primary optional clauses
// function, and finally it will call the primary structured headers function.
//
// Once the parser has completed its work, it will return to the LegalToMarkdown function the final
// contents so that that function may call the appropriate writer for outputting the parsed document
// back to the user.
func LegalToMarkdown(contentsFile string, parametersFile string, outputFile string) {

	contents, parameters := setUp(contentsFile, parametersFile)
	contents, parameters = HandleMixins(contents, parameters)

	headers := SetTheHeaders(contents, parameters)
	contents = HandleTheHeaders(contents, headers)

	writeAFile(outputFile, contents)
}

// MakeYAMLFrontMatter is a convenience function which will parse the contents of a template
// to formulate the YAML Front Matter.
func MakeYAMLFrontMatter(contentsFile string, parametersFile string, outputFile string) {

	contents, parameters := setUp(contentsFile, parametersFile)
	contents = HandleParameterAssembly(contents, parameters)
	writeAFile(outputFile, contents)

}

// MarkdownToPDF renders an input template file to pdf using a webservice hosted on
// https://lmdpdfgen.herokuapp.com/
//
// It runs through the normal parsing system but instead of sending to the standard
// writer it sends the result of the parsing job to the WriteToPdf function which will
// send the result of the parse job to the web service and write the result to the
// output file location.
func MarkdownToPDF(contentsFile string, parametersFile string, outputFile string) {

	contents, parameters := setUp(contentsFile, parametersFile)
	contents, parameters = HandleMixins(contents, parameters)

	headers := SetTheHeaders(contents, parameters)
	contents = HandleTheHeaders(contents, headers)

	WriteToPdf(contents, outputFile)

}

// GetTheParameters is a wrapper function which enables a system to determine what
// the parameters of a given template file are. It will first look at the file to determine
// if there is front matter which can be parsed. If there is front matter that will
// be sent to a function which essentially just marshals the json into a string
// that is returned to any calling function.
//
// If there is no front matter within the template function, the AssembleParametersIntoJSON
// function will be called. This function is analogous to the HandleParameterAssembly function
// which will write to YAML front matter, except that the AseembleParametersIntoJSON function
// will marshall the parameters into a JSON string that will be returned to the calling
// function.
func GetTheParameters(contentsFile string) string {

	_, parameters := setUp(contentsFile, "")

	if len(parameters) == 0 {
		contents, _ := setUp(contentsFile, "")
		return AssembleParametersIntoJSON(contents, parameters)
	} else {
		return jsonizeParameters(parameters)
	}
}

// RawMarkdownToPDF is a function which does not work with written or written_to files
// instead of dealing with files to be read, the function simply parses strings which are
// passed to it programmatically. Otherwise the function logic mirrors MarkdownToPDF.
func RawMarkdownToPDF(rawContents string, rawParameters string) string {

	contents, parameters := setUpRaw(rawContents, rawParameters)
	contents, parameters = HandleMixins(contents, parameters)

	headers := SetTheHeaders(contents, parameters)
	contents = HandleTheHeaders(contents, headers)

	return WriteToPdfRaw(contents)

}

// setUp is a simple convenience function which assists all of the major parsing functions
// in this file. First it will read the contentsFile into memory, along with the included
// partials. Then it will first parse the contentsFile to see if it has front matter. If
// there is front matter these will be unmarshalled into the parameters map.
//
// If a paramaters file is sent to the function, then that file will also be unmarshalled
// and any paramaters which are contained in both the template file and the parameters file
// will be overwritten in favor of the values included in the parameters file.
func setUp(contentsFile string, parametersFile string) (string, map[string]string) {

	// read the template file and integrate any included partials (`@include PARTIAL` within the text)
	contents := ReadAFile(contentsFile)
	contents = importIncludedFiles(contents)

	// once the content files have been read, then move along to parsing the parameters.
	var parameters string
	var amendedParameters map[string]string
	if parametersFile != "" {

		// first pull out of the file, just as we do if there is no specific params file
		var mergedParameters map[string]string
		parameters, contents = parseTemplateToFindParameters(contents)
		mergedParameters = unmarshallParameters(parameters)

		// second read and unmarshall the parameters from the parameters file
		parameters = ReadAFile(parametersFile)
		amendedParameters = unmarshallParameters(parameters)

		// finally, merge the amendedParameters (from the parameters file) into the
		//   mergedParameters (from the content file) such that the amendedParameters
		//   overwritethe mergedParameters.
		amendedParameters = mergeParameters(amendedParameters, mergedParameters)

	} else {

		// if there is no parameters file passed, simply pull the params out of the content file.
		parameters, contents = parseTemplateToFindParameters(contents)
		amendedParameters = unmarshallParameters(parameters)

	}
	return contents, amendedParameters
}

func setUpRaw(contents string, rawParameters string) (string, map[string]string) {

	// once the content files have been read, then move along to parsing the parameters.
	var parameters string
	var amendedParameters map[string]string
	if rawParameters != "" {

		// first pull out of the file, just as we do if there is no specific params file
		var mergedParameters map[string]string
		parameters, contents = parseTemplateToFindParameters(contents)
		mergedParameters = unmarshallParameters(parameters)

		// second read and unmarshall the parameters from the parameters file
		amendedParameters = unmarshallParameters(rawParameters)

		// finally, merge the amendedParameters (from the parameters file) into the
		//   mergedParameters (from the content file) such that the amendedParameters
		//   overwritethe mergedParameters.
		amendedParameters = mergeParameters(amendedParameters, mergedParameters)

	} else {

		// if there is no parameters file passed, simply pull the params out of the content file.
		parameters, contents = parseTemplateToFindParameters(contents)
		amendedParameters = unmarshallParameters(parameters)

	}
	return contents, amendedParameters
}