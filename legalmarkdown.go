package main

import (
    "os"
    "fmt"
    "log"
    "io/ioutil"
    "github.com/codegangsta/cli"
    // "gopkg.in/yaml.v2"
    // "path/filepath"
    // "strings"
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
    contents := read_a_file(c.String("template"))

    var parameters string
    if c.String("parameters") != "" {
        parameters     = read_a_file(c.String("parameters"))
    } else {
        // TODO: parse parameters function goes here.
        parameters     = ""
    }

    // placeholders. TODO: build master LegalToMarkdown(contents string, parameters struct, output string)
    fmt.Print(contents)
    fmt.Print("\n---\n")
    fmt.Print(parameters)
}

func CLIMakeYAMLFrontMatter(c *cli.Context) {
    contents := read_a_file(c.String("template"))
    fmt.Print(contents)
}

func read_a_file(file_to_parse string) (contents string){

    file_buffer, file_read_err := ioutil.ReadFile(file_to_parse)

    if file_read_err != nil {
        log.Fatal(file_read_err)
    }

    contents = string(file_buffer)

    return contents
}