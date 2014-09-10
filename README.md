# Legal Markdown - GoLang

GoLang Port of legal_markdown.

The legal_markdown spec available (sort of) [here](https://github.com/compleatang/legal-markdown/blob/master/README.md) shall apply to this repository and package with one change.

# Modifications to Ruby port of Legal Markdown

* Header information **will not** be kept in the same files as the provisions and other content. So two files (or objects if calling programmatically) should always be passed to the package individually with the header information format passed first and the content file passed second.
* Header information can be passed in `json` or in `yaml` format.
* PDF rendering added as an output format in addition to markdown. PDFs shall be rendered according to the passed template which should conform to the Pandoc PDF template spec.
* Multiple content files can be passed. They will be assembled in sequential order as they are passed.
* There shall be no XML output.
* Only parse 5 levels of headers

The remainder of the spec shall be abided by.

# License

MIT, see LICENSE file.