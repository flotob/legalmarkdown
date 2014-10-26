package main

import (
    "os"
    "fmt"
    "log"
    "github.com/codegangsta/cli"
    // "path/filepath"
)

// manages the performance of the entire parsing job.
//   It accepts a base lmd file to parse along with an optional parameters file.
func main(){
    lmd := cli.NewApp()

    lmd.Name    = "legalmarkdown"
    lmd.Usage   = "Automate your contracting."
    lmd.Version = "0.1.0"
    lmd.Author  = "Eris Industries, Ltd."
    lmd.Email   = "contact@erisindustries.com"

    lmd.Commands = []cli.Command{

      {
        Name:      "tomd",
        ShortName: "m",
        Usage:     "parse to markdown",
        Flags:     []cli.Flag{
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
        Action:    CLILegalToMarkdown,
      },

      {
        Name:      "headers",
        ShortName: "d",
        Usage:     "create yaml frontmatter",
        Flags:     []cli.Flag{
            cli.StringFlag{
              Name:  "t, template",
              Usage: "template file to be parsed",
            },
            cli.StringFlag{
              Name:  "o, output",
              Usage: "output file to be written",
            },
        },
        Action:    CLIMakeYAMLFrontMatter,
      },
    }

    lmd.Run(os.Args)
}

func CLILegalToMarkdown(c *cli.Context) {

    if c.String("template") == "" {
        log.Fatal("Please specify a template file to parse with the --template or -t flag.")
    }

    contents := read_a_file(c.String("template"))
    contents  = import_files(contents)

    // once the content files have been read, then move along to parsing the parameters.
    var parameters         string
    var amended_parameters map[string]string
    if c.String("parameters") != "" {
        // first pull out of the file, just as we do if there is no specific params file
        var merged_parameters map[string]string
        parameters, contents = parse_template_to_find_parameters(contents)
        merged_parameters    = unmarshall_parameters(parameters)
        // second read and unmarshall the parameters from the parameters file
        parameters           = read_a_file(c.String("parameters"))
        amended_parameters   = unmarshall_parameters(parameters)
        // finally, merge the amended_parameters (from the parameters file) into the
        //   merged_parameters (from the content file) such that the amended_parameters
        //   overwritethe merged_parameters.
        amended_parameters   = merge_parameters(amended_parameters, merged_parameters)
    } else {
        // if there is no parameters file passed, simply pull the params out of the content file.
        parameters, contents = parse_template_to_find_parameters(contents)
        amended_parameters   = unmarshall_parameters(parameters)
    }

    // placeholders. TODO: build master LegalToMarkdown(contents string, parameters struct, output string)
    fmt.Print(contents)
    fmt.Print("\n--*****--\n")
    fmt.Print(amended_parameters)
}

func CLIMakeYAMLFrontMatter(c *cli.Context) {
    contents := read_a_file(c.String("template"))
    fmt.Print(contents)
}