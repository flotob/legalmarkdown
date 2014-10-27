// The LegalMarkdown package contains the parser functions to automate and simplify contracting.
// You can easily transform your contracting by using the functionality of how programmers work
// into your legal practice. By integrating such programming concepts as partials, templates,
// boolean triggers, and the like (which sound fancy but are very simple to explain), your
// transactional practice will be greatly simplified.
package legalmarkdown

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
	// "path/filepath"
)

// main parses the command line inputs and routes the commands to the appropriate wrapper
// functions.
func main() {
	lmd := cli.NewApp()

	lmd.Name = "legalmarkdown"
	lmd.Usage = "Automate your contracting."
	lmd.Version = "0.1.0"
	lmd.Author = "Eris Industries, Ltd."
	lmd.Email = "contact@erisindustries.com"

	lmd.Commands = []cli.Command{

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

	lmd.Run(os.Args)
}

func cliLegalToMarkdown(c *cli.Context) {

	if c.String("template") == "" {
		log.Fatal("Please specify a template file to parse with the --template or -t flag.")
	}

	contents   := c.String("template")
	parameters := c.String("parameters")
    output     := c.String("output")

	LegalToMarkdown(contents, parameters, output)
}

func cliMakeYAMLFrontMatter(c *cli.Context) {

	contents := read_a_file(c.String("template"))
	fmt.Print(contents)

}

// LegalToMarkdown is the primary function which controls parsing a template document into a markdown
// result when the parsing library is called from the command line. Two strings which are filenames
// should be passed to the function. The parameters string may be an empty string. The function first
// parses and reads the command sent from the command line, and then reads the template file. After
// this, the function pulls into the template file any partials which have been included into the
// template with the
//   @include {{PARTIAL}}
// flag within the text of the primary template file that has been called.
//
// Then the function reads the paramaters from a parameters file, the template file, or both. In the
// case where the parameters are contained in both a parameters file and in the template file, the
// parameters in the template file are considered as defaults which are overridden by parameters passed
// to the function via the paramaters file.
//
// Finally once the function has prepared the
//   contents
// and
//   parameters
// from the various passed files and built a cohesive set of
//   contents
// and
//   parameters
//
// These are passed to the primary entrance function to the parsing process.
func LegalToMarkdown(contents string, parameters string, output string) {

	// read the template file and integrate any included partials (`@include PARTIAL` within the text)
	contents = read_a_file(contents)
	contents = import_files(contents)

	// once the content files have been read, then move along to parsing the parameters.
	var amended_parameters map[string]string
	if parameters != "" {

		// first pull out of the file, just as we do if there is no specific params file
		var merged_parameters map[string]string
		parameters, contents = parse_template_to_find_parameters(contents)
		merged_parameters = unmarshall_parameters(parameters)

		// second read and unmarshall the parameters from the parameters file
		parameters = read_a_file(parameters)
		amended_parameters = unmarshall_parameters(parameters)

		// finally, merge the amended_parameters (from the parameters file) into the
		//   merged_parameters (from the content file) such that the amended_parameters
		//   overwritethe merged_parameters.
		amended_parameters = merge_parameters(amended_parameters, merged_parameters)

	} else {

		// if there is no parameters file passed, simply pull the params out of the content file.
		parameters, contents = parse_template_to_find_parameters(contents)
		amended_parameters = unmarshall_parameters(parameters)

	}

	contents = legalToMarkdownParser(contents, amended_parameters)

    // TODO: call the appropriate writer
    write_a_file(output, contents)
}

// MakeYAMLFrontMatter is a convenience function which will parse the contents of a template
// to formulate the YAML Front Matter.
func MakeYAMLFrontMatter(contents string) (string) {
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


	return contents
}
