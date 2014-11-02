// The LegalMarkdown package contains the parser functions to automate and simplify contracting.
// You can easily transform your contracting by using the functionality of how programmers work
// into your legal practice. By integrating such programming concepts as partials, templates,
// boolean triggers, and the like (which sound fancy but are very simple to explain), your
// transactional practice will be greatly simplified.
package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/eris-ltd/legalmarkdown/lmd"
	"log"
	"os"
)

// main parses the command line inputs and routes the commands to the appropriate wrapper
// functions.
func main() {
	legalmd := cli.NewApp()

	legalmd.Name = "legalmarkdown"
	legalmd.Usage = "Automate your contracting."
	legalmd.Version = "0.1.0"
	legalmd.Author = "Eris Industries, Ltd."
	legalmd.Email = "contact@erisindustries.com"

	legalmd.Commands = []cli.Command{

		{
			Name:      "tomd",
			ShortName: "m",
			Usage:     "parse to markdown",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "t, template",
					Usage: "template file to be parsed",
				},
				cli.StringFlag{
					Name:  "p, parameters",
					Usage: "parameters file to be parsed",
				},
				cli.StringFlag{
					Name:  "o, output",
					Usage: "output file to be written",
				},
			},
			Action: cliLegalToMarkdown,
		},

		{
			Name:      "headers",
			ShortName: "d",
			Usage:     "create yaml frontmatter",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "t, template",
					Usage: "template file to be parsed",
				},
				cli.StringFlag{
					Name:  "o, output",
					Usage: "output file to be written",
				},
			},
			Action: cliMakeYAMLFrontMatter,
		},
	}

	legalmd.Run(os.Args)
}

func cliLegalToMarkdown(c *cli.Context) {

	if c.String("template") == "" {
		log.Fatal("Please specify a template file to parse with the --template or -t flag.")
	}

	if c.String("output") == "" {
		log.Fatal("Please specify an output file to write to with the --output or -o flag.")
	}

	contents := c.String("template")
	parameters := c.String("parameters")
	output := c.String("output")

	LegalToMarkdown(contents, parameters, output)
}

func cliMakeYAMLFrontMatter(c *cli.Context) {

	contents := lmd.ReadAFile(c.String("template"))
	fmt.Print(contents)

}

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
// These are passed to the primary entrance function to the parsing process.
func LegalToMarkdown(contents string, parameters_file string, output string) {

	// read the template file and integrate any included partials (`@include PARTIAL` within the text)
	contents = lmd.ReadAFile(contents)
	contents = lmd.ImportIncludedFiles(contents)

	// once the content files have been read, then move along to parsing the parameters.
	var parameters string
	var amended_parameters map[string]string
	if parameters_file != "" {

		// first pull out of the file, just as we do if there is no specific params file
		var merged_parameters map[string]string
		parameters, contents = lmd.ParseTemplateToFindParameters(contents)
		merged_parameters = lmd.UnmarshallParameters(parameters)

		// second read and unmarshall the parameters from the parameters file
		parameters = lmd.ReadAFile(parameters_file)
		amended_parameters = lmd.UnmarshallParameters(parameters)

		// finally, merge the amended_parameters (from the parameters file) into the
		//   merged_parameters (from the content file) such that the amended_parameters
		//   overwritethe merged_parameters.
		amended_parameters = lmd.MergeParameters(amended_parameters, merged_parameters)

	} else {

		// if there is no parameters file passed, simply pull the params out of the content file.
		parameters, contents = lmd.ParseTemplateToFindParameters(contents)
		amended_parameters = lmd.UnmarshallParameters(parameters)

	}

	contents = legalToMarkdownParser(contents, amended_parameters)

	lmd.WriteAFile(output, contents)
}

// MakeYAMLFrontMatter is a convenience function which will parse the contents of a template
// to formulate the YAML Front Matter.
func MakeYAMLFrontMatter(contents string) string {
	// TODO: it all.
	return contents
}

// legalToMarkdownParser is the overseer of the parsing functionality. The contents of the file
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
func legalToMarkdownParser(contents string, parameters map[string]string) string {
	contents, parameters = lmd.HandleMixins(contents, parameters)
	headers := lmd.SetTheHeaders(contents, parameters)
	contents = lmd.HandleTheHeaders(contents, headers)
	return contents
}
