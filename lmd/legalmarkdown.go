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

// MarkdownToPDF
func MarkdownToPDF(contentsFile string, parametersFile string, outputFile string) {

	contents, parameters := setUp(contentsFile, parametersFile)
	contents, parameters = HandleMixins(contents, parameters)

	headers := SetTheHeaders(contents, parameters)
	contents = HandleTheHeaders(contents, headers)

	writeToPdf(contents, outputFile)

}

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
