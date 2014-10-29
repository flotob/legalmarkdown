Master [![Build Status](https://travis-ci.org/eris-ltd/legalmarkdown.svg?branch=master)](https://travis-ci.org/eris-ltd/legalmarkdown) || Develop [![Build Status](https://travis-ci.org/eris-ltd/legalmarkdown.svg?branch=develop)](https://travis-ci.org/eris-ltd/legalmarkdown)

## Legal Markdown - GoLang

Go port of legalmarkdown. The legalmarkdown spec is available [here](https://github.com/compleatang/legal-markdown/blob/master/README.md). That specification shall apply to this repository and package with a few changes noted below.

## Modifications from Ruby version of Legal Markdown

* Header information may be kept in the same or in different files as the provisions and other content. So two files (or objects if calling programmatically) may be passed to the package individually with the header information format passed first and the content file passed second.
* Header information can be passed in `json` or in `yaml` format.
* PDF rendering added as an output format in addition to markdown.
* Multiple content files can be passed. They will be assembled in sequential order as they are passed.
* There shall be no XML output.

The remainder of the spec shall be abided by.

## License

MIT, see LICENSE file.